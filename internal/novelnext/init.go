package novelnext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"novel-scraper-bot/internal/scrapper"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type NovelNextScrapper struct {
	StartingChapterNumber int
	EndingChapterNumber   int
	HomePageURL           string
	ChapterURls           []string
	Content               []scrapper.Chapter
}

type response struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
	HTML    string `json:"html"`
}

func CreateNovelNextScrapper(url string, start int, end int) scrapper.Scrapper {
	return &NovelNextScrapper{
		HomePageURL:           url,
		StartingChapterNumber: start,
		EndingChapterNumber:   end,
		Content:               make([]scrapper.Chapter, end-(start-1)),
	}
}

func (n *NovelNextScrapper) FetchAllLinksOfChapters() error {
	slog.Info("fetching home page html content")
	html, err := fetchHomepageHTML(n.HomePageURL)
	if err != nil {
		return err
	}
	slog.Info("completed fetching home page")

	pageReader := strings.NewReader(html)
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
	return nil
}

func (n *NovelNextScrapper) FetchAllChaptersContent() error {
	sem := semaphore.NewWeighted(10)
	g := new(errgroup.Group)

	slog.Info("fetching chapters content")
	for i, url := range n.ChapterURls[n.StartingChapterNumber-1 : n.EndingChapterNumber] {
		g.Go(func() error {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				return fmt.Errorf("error in acquiring semaphore, url: %v, counter: %v, error: %v", url, i, err)
			}
			defer sem.Release(1)

			fmt.Println("fetch url", i, url)
			html, err := fetchHTML(url)
			if err != nil {
				return err
			}
			fmt.Println("fetch finished url", i, url)

			pageReader := strings.NewReader(html)
			doc, err := goquery.NewDocumentFromReader(pageReader)
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

			n.Content[i] = scrapper.Chapter{Title: title, Content: content}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (n *NovelNextScrapper) GetAllChaptersContent() []scrapper.Chapter {
	return n.Content
}

func fetchHTML(url string) (string, error) {
	scriptPath := filepath.Join("internal", "novelnext", "scripts", "get-html.py")
	cmd := exec.Command("python3", scriptPath, "--url", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	// run the python script
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	// read the response
	var resp response
	if err = json.Unmarshal(out.Bytes(), &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("error in fetching home page, error=%s", resp.Reason)
	}
	return resp.HTML, nil
}

func fetchHomepageHTML(url string) (string, error) {
	scriptPath := filepath.Join("internal", "novelnext", "scripts", "get-homepage-html.py")
	cmd := exec.Command("python3", scriptPath, "--url", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	// run the python script
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	// read the response
	var resp response
	if err = json.Unmarshal(out.Bytes(), &resp); err != nil {
		return "", err
	}
	if !resp.Success {
		return "", fmt.Errorf("error in fetching home page, error=%s", resp.Reason)
	}
	return resp.HTML, nil
}
