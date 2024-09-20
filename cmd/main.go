package main

import (
	"log/slog"
	"novel-scraper-bot/interanl/app"
)

func main() {
	if err := app.MakeEpub("", 1, 500, ""); err != nil {
		slog.Error("error in making epub", "error", err)
	}
}
