package app

import (
	"novel-scraper-bot/internal/novelnext"
)

func MakeEpub(url string, startChapter float64, endChapter float64, filename string) error {
	novel := novelnext.Init(url, startChapter, endChapter)
	if err := novel.FetchAllChaptersContent(); err != nil {
		return err
	}
	return nil
}
