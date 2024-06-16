package scrappers

type Scrapper interface {
	FetchAllLinksOfChapters() error
	FetchAllChaptersContent() error
	GetAllChaptersContent() map[int]Chapter
}

type Chapter struct {
	Title   string
	Content string
}
