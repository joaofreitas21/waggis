package views

type CtfRender struct {
	Name   string
	Active bool
}

type CtfCard struct {
	Slug          string
	Title         string
	Platform      string
	Difficulty    string
	DateFormatted string
	Tags          []CtfRender
	Description   string
}

type CtfMetadata struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Platform    string   `json:"platform"`
	Difficulty  string   `json:"difficulty"`
	Date        string   `json:"date"`
	Tags        []string `json:"tags"`
	Complete    bool     `json:"complete"`
	Description string   `json:"description"`
}

