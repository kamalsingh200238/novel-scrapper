package scrapper

type Scrapper interface {
	FetchAllLinksOfChapters() error
	FetchAllChaptersContent() error
	GetAllChaptersContent() []Chapter
}

type Chapter struct {
	Number  float64
	Title   string
	Content string
}
