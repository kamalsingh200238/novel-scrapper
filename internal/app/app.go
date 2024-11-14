package app

import (
	"fmt"
	"novel-scraper-bot/internal/novelnext"
)

func MakeEpub(url string, startChapter float64, endChapter float64, filename string) error {
	novel := novelnext.Init(url, startChapter, endChapter)
	fmt.Println(novel)
	return fmt.Errorf("novel next does not work anymore")
}
