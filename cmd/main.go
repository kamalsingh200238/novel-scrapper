package main

import (
	"flag"
	"log/slog"
	"novel-scraper-bot/internal/app"
)

func main() {
	url := flag.String("url", "", "URL of the novel")
	startChapter := flag.Float64("start", 0, "starting chapter number")
	endChapter := flag.Float64("end", 0, "ending chapter number")
	fileName := flag.String("filename", "", "output filename for the epub")

	flag.Parse()

	if err := app.MakeEpub(*url, *startChapter, *endChapter, *fileName); err != nil {
		slog.Error("error in making epub",
			"error", err,
			"url", *url,
			"file name", *fileName,
			"start chatper", *startChapter,
			"end chatper", *endChapter,
		)
	}
}
