package scrapper

type Scrapper interface {
	FetchAllLinksOfChapters() error
	FetchAllChaptersContent() error
	GetAllChaptersContent() []Chapter
}

type Chapter struct {
	Title   string
	Content string
}

