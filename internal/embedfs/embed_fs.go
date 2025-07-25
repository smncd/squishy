package embedfs

import (
	"embed"
)

//go:embed *
var FS embed.FS

type ErrorPageData struct {
	Title       string
	Description string
	Error       string
}
