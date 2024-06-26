package scrappers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type NovelNextScrapper struct {
	StartingChapterNumber int
	EndingChapterNumber   int
	HomePageURL           string
	ChapterURls           []string
	Content               []Chapter
}

func CreateNovelNextScrapper(url string, start int, end int) Scrapper {
	return &NovelNextScrapper{
		HomePageURL:           url,
		StartingChapterNumber: start,
		EndingChapterNumber:   end,
		Content:               make([]Chapter, end-(start-1)),
	}
}

func (n *NovelNextScrapper) FetchAllLinksOfChapters() error {
	// get the html content of the home page
	// reason for using chrome dp: novel next uses javascript to load all the urls for the page
	fmt.Println("fetching home page html content")
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var htmlContent string
	// immitate infinte scroll
	if err := chromedp.Run(ctx,
		chromedp.Navigate(n.HomePageURL),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.KeyEvent(kb.End),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.OuterHTML("html", &htmlContent),
	); err != nil {
		return fmt.Errorf("error in fetching chapter urls from the home page, error: %v", err)
	}
	fmt.Println("completed fetching home page")

	fmt.Println("starting to parse chapter links from the html content")
	pageReader := strings.NewReader(htmlContent)
	doc, err := goquery.NewDocumentFromReader(pageReader)
	if err != nil {
		return fmt.Errorf("error in parsing home page html, error: %v", err)
	}

	ulElements := doc.Find("ul.list-chapter")
	ulElements.Each(func(i int, ul *goquery.Selection) {
		ul.Find("a").Each(func(j int, a *goquery.Selection) {
			// Get the href attribute of the current a element
			href, exists := a.Attr("href")
			if exists {
				n.ChapterURls = append(n.ChapterURls, href)
			}
		})
	})
	fmt.Println("finished parsing chapter links from the html content")

	return nil
}

func (n *NovelNextScrapper) FetchAllChaptersContent() error {
	sem := semaphore.NewWeighted(10)
	g := new(errgroup.Group)

	fmt.Println("fetching chapter content")
	for i, url := range n.ChapterURls[n.StartingChapterNumber-1 : n.EndingChapterNumber] {
		i, url := i, url
		g.Go(func() error {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				return fmt.Errorf("error in acquiring semaphore, url: %v, counter: %v, error: %v", url, i, err)
			}
			defer sem.Release(1)

			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("error in fetching the content, url: %v, counter: %v error: %v", url, i, err)
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return fmt.Errorf("error in parsing the content of page, url: %v, counter: %v, error: %v", url, i, err)
			}

			title := strings.TrimSpace(doc.Find(".chr-title").First().Text())

			var content string

			var processNode func(*goquery.Selection)
			processNode = func(s *goquery.Selection) {
				if s.Children().Length() == 0 {
					tagName := goquery.NodeName(s)
					if tagName != "script" {
						content += fmt.Sprintf("%v\n", strings.TrimSpace(s.Text()))
					}
				} else {
					s.Children().Each(func(i int, selection *goquery.Selection) {
						processNode(selection)
					})
				}
			}

			doc.Find("#chr-content").Children().Each(func(i int, s *goquery.Selection) {
				processNode(s)
			})

			n.Content[i] = Chapter{Title: title, Content: content}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (n *NovelNextScrapper) GetAllChaptersContent() []Chapter {
	return n.Content
}
