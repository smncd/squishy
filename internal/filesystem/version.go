package filesystem

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:generate sh -c "cp ../../.version ./_version"

//go:embed _version
var staticFS embed.FS

var VersionString string

func LoadVersionString() error {
	data, err := fs.ReadFile(staticFS, "_version")
	if err != nil {
		return fmt.Errorf("could not read embedded _version file: %w", err)
	}

	VersionString = string(data)

	return nil
}
