package main

import (
	"fmt"
	"log/slog"
	"novel-scraper-bot/interanl/scrappers"

	"github.com/bmaupin/go-epub"
)

func main() {
	url := "https://novelnextz.com/novelnextz/dragon-marked-war-god#tab-chapters-title"
	start := 650
	end := 700
	novel := scrappers.CreateNovelNextScrapper(url, start, end)

	err := novel.FetchAllLinksOfChapters()
	if err != nil {
		slog.Error("in fetching all urls",
			"error", err,
		)
	}

	err = novel.FetchAllChaptersContent()
	if err != nil {
		slog.Error("in fetching content for all chapters",
			"error", err,
		)
	}

	m := novel.GetAllChaptersContent()
	fileName := fmt.Sprintf("dragon_marked_war_god_%v_%v.epub", start, end)
	novelHeader := fmt.Sprintf("dragon_marked_war_god_%v_%v", start, end)
	e := epub.NewEpub(novelHeader)

	fmt.Println("Writing in file")
	for _, v := range m {
		e.AddSection(v.Content, v.Title, "", "")
	}
	err = e.Write(fileName)
	if err != nil {
		slog.Error("in writing in the epub file",
			"error", err,
		)
	}
	fmt.Println("finished writing in file")
}
