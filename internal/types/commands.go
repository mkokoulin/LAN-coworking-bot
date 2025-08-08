package types

type Event struct {
	Capacity     string
	Date         string
	Description  string
	ExternalLink string
	Id           string
	Img          string
	Link         string
	Name         string
	ShowForm     bool
	Type         string
}

type LangOption struct {
	Code string // "en" / "ru"
	Label string // "🇺🇸 English"
}