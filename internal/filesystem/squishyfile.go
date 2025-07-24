package filesystem

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type meta struct {
	modifiedTime time.Time
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

	FilePath string         `validate:"required"`
	Config   config         `yaml:"config" json:"config"`
	Routes   map[string]any `yaml:"routes" json:"routes" validate:"required"`
}

// Loads SquishyFile from filesystem
func (s *SquishyFile) Load() error {
	s.Config = config{
		Debug: false,
		Host:  "localhost",
		Port:  "1394",
	}

	err := loadFromFile(s.FilePath, &s)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(s)
	if err != nil {
		return err
	}

	err = s.StoreFileModTime()
	if err != nil {
		return err
	}

	return nil
}

func (s *SquishyFile) GetFileModTime() (*time.Time, error) {
	fileInfo, err := os.Stat(s.FilePath)
	if err != nil {
		return nil, err
	}

	modTime := fileInfo.ModTime()

	return &modTime, nil
}

func (s *SquishyFile) StoreFileModTime() error {
	modTime, err := s.GetFileModTime()
	if err != nil {
		return err
	}

	s.meta.modifiedTime = *modTime

	return nil
}

func (s *SquishyFile) FileUpdatedSinceLastLoad() (bool, error) {
	modTime, err := s.GetFileModTime()
	if err != nil {
		return false, err
	}

	return modTime.Compare(s.meta.modifiedTime) != 0, nil
}

func (s *SquishyFile) RefetchRoutes() error {
	updated, err := s.FileUpdatedSinceLastLoad()
	if err != nil {
		return err
	}

	if updated {
		var newData SquishyFile
		log.Println("squishyfile has new mod time, loading again...")

		err := loadFromFile(s.FilePath, &newData)
		if err != nil {
			return err
		}

		s.Routes = newData.Routes

		err = s.StoreFileModTime()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SquishyFile) LookupRoutePath(path string) (string, bool) {
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
			return "", false
		}

		result, ok = currentLevel[key]
		if !ok {
			return "", false
		}

		if i == len(keys)-1 {
			if level, ok := result.(map[string]any); ok {
				result, ok = level[indexKey]
				if !ok {
					return "", false
				}
			}

			break
		}
	}

	reply, ok := result.(string)
	if !ok {
		return "", false
	}

	return reply, true

}

// Loads SquishyFile from filesystem.
func loadFromFile(filePath string, data any) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		return err
	}

	return nil
}
