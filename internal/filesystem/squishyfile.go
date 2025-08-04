package filesystem

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type meta struct {
	filePath     string    `validate:"required"`
	modifiedTime time.Time `validate:"required"`
	logger       *log.Logger
}

type config struct {
	Debug bool   `yaml:"debug" json:"debug"`
	Host  string `yaml:"host" json:"host" validate:"required"`
	Port  string `yaml:"port" json:"port" validate:"required"`
}

// The SquishyFile holds the main configuration data for the project,
// including settings as well as all the available routes
type SquishyFile struct {
	meta meta

	Config config         `yaml:"config" json:"config"`
	Routes map[string]any `yaml:"routes" json:"routes" validate:"required"`
}

func (s *SquishyFile) SetFilePath(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	s.meta.filePath = filePath

	return true
}

func (s *SquishyFile) SetLogger(logger *log.Logger) {
	s.meta.logger = logger
}

// Loads SquishyFile from filesystem
func (s *SquishyFile) Load() error {
	s.Config = config{
		Debug: false,
		Host:  "localhost",
		Port:  "1394",
	}

	if s.meta.logger == nil {
		s.meta.logger = log.New(os.Stderr, "", 0)
	}

	err := loadFromFile(s.meta.filePath, &s)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	err = s.StoreFileModTime()
	if err != nil {
		return fmt.Errorf("failed to store file modification time: %w", err)
	}

	return nil
}

func (s *SquishyFile) GetFileModTime() (*time.Time, error) {
	fileInfo, err := os.Stat(s.meta.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	modTime := fileInfo.ModTime()

	return &modTime, nil
}

func (s *SquishyFile) StoreFileModTime() error {
	modTime, err := s.GetFileModTime()
	if err != nil {
		return fmt.Errorf("failed to get file modification time: %w", err)
	}

	s.meta.modifiedTime = *modTime

	return nil
}

func (s *SquishyFile) FileUpdatedSinceLastLoad() (bool, error) {
	modTime, err := s.GetFileModTime()
	if err != nil {
		return false, fmt.Errorf("failed to get file modification time: %w", err)
	}

	return modTime.Compare(s.meta.modifiedTime) != 0, nil
}

func (s *SquishyFile) RefetchRoutes() error {
	updated, err := s.FileUpdatedSinceLastLoad()
	if err != nil {
		return fmt.Errorf("failed to check if file was updated since last load: %w", err)
	}

	if updated {
		var newData SquishyFile
		s.meta.logger.Println("squishyfile has new mod time, loading again...")

		err := loadFromFile(s.meta.filePath, &newData)
		if err != nil {
			return fmt.Errorf("failed to get file modification time: %w", err)
		}

		s.Routes = newData.Routes

		err = s.StoreFileModTime()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SquishyFile) LookupRouteUrlFromPath(path string) (string, error) {
	indexKey := "_index"
	var keys []string

	for part := range strings.SplitSeq(strings.ReplaceAll(path, ":", "/"), "/") {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart == "" {
			continue
		}
		keys = append(keys, trimmedPart)
	}

	if len(keys) == 0 {
		keys = append(keys, indexKey)
	}

	var result any = s.Routes

	for i, key := range keys {
		currentLevel, ok := result.(map[string]any)
		if !ok {
			return "", errors.New("result is not of type map[string]any")
		}

		result, ok = currentLevel[key]
		if !ok {
			return "", fmt.Errorf("key %s not found", key)
		}

		if i == len(keys)-1 {
			if level, ok := result.(map[string]any); ok {
				result, ok = level[indexKey]
				if !ok {
					return "", fmt.Errorf("key %s not found", indexKey)
				}
			}

			break
		}
	}

	reply, ok := result.(string)
	if !ok {
		return "", errors.New("result is not string")
	}

	return reply, nil

}

// Loads SquishyFile from filesystem.
func loadFromFile(filePath string, data any) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file (%s): %w", filePath, err)
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}
