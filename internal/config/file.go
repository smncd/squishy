package config

import (
	"fmt"
	"os"
	"time"
)

type file struct {
	Path         string `validate:"required"`
	modifiedTime time.Time
}

func (f *file) GetModTime() (*time.Time, error) {
	fileInfo, err := os.Stat(f.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	modTime := fileInfo.ModTime()

	return &modTime, nil
}

func (f *file) StoreModTime() error {
	modTime, err := f.GetModTime()
	if err != nil {
		return fmt.Errorf("failed to get file modification time: %w", err)
	}

	f.modifiedTime = *modTime

	return nil
}

func (f *file) UpdatedSinceLastLoad() (bool, error) {
	modTime, err := f.GetModTime()
	if err != nil {
		return false, fmt.Errorf("failed to get file modification time: %w", err)
	}

	return modTime.Compare(f.modifiedTime) != 0, nil
}
