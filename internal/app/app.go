package app

import (
	"fmt"
	"log/slog"
	"novel-scraper-bot/internal/novelnext"

	"github.com/bmaupin/go-epub"
)

func MakeEpub(url string, startChapter int, endChapter int, filename string) error {
	novel := novelnext.CreateNovelNextScrapper(url, startChapter, endChapter)
	err := novel.FetchAllLinksOfChapters()
	if err != nil {
		return err
	}

	err = novel.FetchAllChaptersContent()
	if err != nil {
		return err
	}

	m := novel.GetAllChaptersContent()
	fmt.Println(m)

	fileName := fmt.Sprintf("%s_%d_%d.epub", filename, startChapter, endChapter)
	novelHeader := fmt.Sprintf("%s_%d_%d", filename, startChapter, endChapter)
	e := epub.NewEpub(novelHeader)

	slog.Info("writing in file")
	for _, v := range m {
		e.AddSection(v.Content, v.Title, "", "")
	}
	err = e.Write(fileName)
	if err != nil {
	  return err
	}
	slog.Info("finished writing in file")

	return nil
}
