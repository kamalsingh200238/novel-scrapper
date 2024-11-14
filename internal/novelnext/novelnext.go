package novelnext

import "novel-scraper-bot/internal/scrapper"

type NovelNext struct {
	startChapter float64
	endChapter   float64
	homepageURL  string
	chapters     []scrapper.Chapter
}

func Init(url string, start float64, end float64) scrapper.Scrapper {
	return &NovelNext{
		homepageURL:  url,
		startChapter: start,
		endChapter:   end,
	}
}

func (n *NovelNext) FetchAllLinksOfChapters() error            { return nil }
func (n *NovelNext) FetchAllChaptersContent() error            { return nil }
func (n *NovelNext) GetAllChaptersContent() []scrapper.Chapter { return []scrapper.Chapter{} }
