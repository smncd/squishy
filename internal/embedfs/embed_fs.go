package embedfs

import (
	"embed"
)

//go:embed *
var FS embed.FS
