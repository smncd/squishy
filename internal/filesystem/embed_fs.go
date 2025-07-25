package filesystem

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:generate sh -c "cp ../../.version ./embed/_version"

//go:embed embed/*
var EmbedFS embed.FS

var VersionString string

func LoadVersionString() error {
	data, err := fs.ReadFile(EmbedFS, "embed/_version")
	if err != nil {
		return fmt.Errorf("could not read embedded _version file: %w", err)
	}

	VersionString = string(data)

	return nil
}
