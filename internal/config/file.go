package config

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type file struct {
	path         string `validate:"required"`
	modifiedTime time.Time
}

func (f *file) Load(out any) error {
	rawData, err := os.ReadFile(f.path)
	if err != nil {
		return fmt.Errorf("failed to read config file (%s): %w", f.path, err)
	}

	if err := yaml.Unmarshal(rawData, out); err != nil {
		return fmt.Errorf("failed to unmarshal YAML config: %w", err)
	}

	val := reflect.ValueOf(out)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("config pointer is nil")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a struct or a pointer to a struct, is: %v", val.Kind())
	}

	validate := validator.New()
	err = validate.Struct(out)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	err = f.StoreModTime()
	if err != nil {
		return fmt.Errorf("failed to store file modification time: %w", err)
	}

	return nil
}

func (f *file) GetModTime() (*time.Time, error) {
	fileInfo, err := os.Stat(f.path)
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
