package scrapper

type Scrapper interface {
	FetchAllChaptersContent() error
	GetAllChaptersContent() []Chapter
}

type Chapter struct {
	Number  float64 `json:"number"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
}
