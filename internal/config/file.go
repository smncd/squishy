package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type file struct {
	Path         string `arg:"-C,--config-file,env:CONFIG_FILE" default:"./squishy.yaml"`
	modifiedTime time.Time
}

func (f *file) Load(out any) error {
	ext := filepath.Ext(f.Path)

	supportedExts := map[string]bool{
		".yaml": true,
		".yml":  true,
		".json": true,
		".toml": true,
	}

	if !supportedExts[ext] {
		return fmt.Errorf("unsupported config file extension: %s", ext)
	}

	rawData, err := os.ReadFile(f.Path)
	if err != nil {
		return fmt.Errorf("failed to read config file (%s): %w", f.Path, err)
	}

	switch ext := filepath.Ext(f.Path); ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(rawData, out); err != nil {
			return fmt.Errorf("failed to unmarshal YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(rawData, out); err != nil {
			return fmt.Errorf("failed to unmarshal JSON config: %w", err)
		}
	case ".toml":
		if _, err := toml.Decode(string(rawData), out); err != nil {
			return fmt.Errorf("failed to unmarshal TOML config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file extension: %s", ext)
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
