package main

import (
	"log/slog"
	"novel-scraper-bot/internal/app"
)

func main() {
	if err := app.MakeEpub("https://novel-next.com/novel/supreme-magus-novel#tab-chapters-title", 1, 50, "temp"); err != nil {
		slog.Error("error in making epub", "error", err)
	}
}
