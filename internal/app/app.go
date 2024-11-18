package app

import (
	"novel-scraper-bot/internal/novelnext"
	"novel-scraper-bot/internal/scrapper"
)

func MakeEpub(url string, startChapter, endChapter float64, filename string) error {
	novel := novelnext.Init(url, startChapter, endChapter)
	if err := novel.FetchAllChaptersContent(); err != nil {
		return err
	}
	if err := scrapper.MakeEpub(novel.GetAllChaptersContent(), filename, startChapter, endChapter); err != nil {
		return err
	}
	return nil
}
