package templates

import "embed"

//go:embed *.html
var FS embed.FS

type ErrorPageData struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Error       string
}
