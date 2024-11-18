package novelnext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"novel-scraper-bot/internal/scrapper"
	"os/exec"
	"strings"
)

type NovelNext struct {
	startChapter float64
	endChapter   float64
	homepageURL  string
	chapters     []scrapper.Chapter
}

type ScriptData struct {
	Data  []scrapper.Chapter `json:"data"`
	Error string             `json:"error"`
}

func Init(url string, start, end float64) scrapper.Scrapper {
	return &NovelNext{
		homepageURL:  url,
		startChapter: start,
		endChapter:   end,
		chapters:     []scrapper.Chapter{},
	}
}

func (n *NovelNext) FetchAllChaptersContent() error {
	script := "scripts/novelnext/fetch-content.py"
	cmd := exec.Command(
		".venv/bin/python3",
		script,
		"--url", n.homepageURL,
		"--start", fmt.Sprintf("%f", n.startChapter),
		"--end", fmt.Sprintf("%f", n.endChapter),
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error in executing python script, error=%v", err)
	}

	output := out.String()
	first := strings.Index(output, "{")
	last := strings.LastIndex(output, "}")
	if first == -1 || last == -1 {
		return fmt.Errorf("no json returned from python script")
	}
	jsonString := output[first : last+1]

	var data ScriptData
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		fmt.Println(err)
	}

	if data.Error != "" {
		return fmt.Errorf("error while executing python script, error=%v", data.Error)
	}

	n.chapters = data.Data

	return nil
}

func (n *NovelNext) GetAllChaptersContent() []scrapper.Chapter { return n.chapters }
