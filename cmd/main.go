package main

import (
	"fmt"
	"log/slog"
	"novel-scraper-bot/interanl/scrappers"
)

func main() {
	url := "https://novelnextz.com/novelnextz/dragon-marked-war-god#tab-chapters-title"
	start := 1
	end := 50
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
	for i := range end - (start - 1) {
		fmt.Println(i, m[i])
	}
}
