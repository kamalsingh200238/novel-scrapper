package scrapper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/go-epub"
)

type Scrapper interface {
	FetchAllChaptersContent() error
	GetAllChaptersContent() []Chapter
}

type Chapter struct {
	Number  float64 `json:"number"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
}

func MakeEpub(chapters []Chapter, fileName string, start, end float64) error {
	// Format the book title and output file path
	bookTitle := fmt.Sprintf("%s_%0.2f_%0.2f", strings.ReplaceAll(fileName, " ", "_"), start, end)
	outputDir := "./output"
	bookName := filepath.Join(outputDir, fmt.Sprintf("%s.epub", bookTitle))

	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("cannot create output directory, error=%v", err)
	}

	// Create a new EPUB
	e, err := epub.NewEpub(bookTitle)
	if err != nil {
		return fmt.Errorf("cannot create epub, error=%v", err)
	}

	// Set metadata
	e.SetAuthor("Scrapped Novel")

	// Add chapters as sections
	for _, chapter := range chapters {
		content := fmt.Sprintf("<h1>%s</h1><p>%s</p>", chapter.Title, chapter.Content)
		_, err = e.AddSection(content, chapter.Title, "", "")
		if err != nil {
			return fmt.Errorf("cannot add chapter '%s', error=%v", chapter.Title, err)
		}
	}

	// Write the EPUB to file
	err = e.Write(bookName)
	if err != nil {
		return fmt.Errorf("cannot write epub file, error=%v", err)
	}

	return nil
}
