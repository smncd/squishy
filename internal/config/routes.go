package config

import (
	"errors"
	"fmt"
	"strings"
)

type Routes struct {
	file   file
	Routes map[string]any `yaml:"routes" json:"routes" validate:"required"`
}

func (r *Routes) Load() error {
	err := r.file.Load(r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Routes) Refetch() error {
	updated, err := r.file.UpdatedSinceLastLoad()
	if err != nil {
		return fmt.Errorf("failed to check if file was updated since last load: %w", err)
	}

	if updated {
		var newData Routes
		// TODO: reimplement logger
		// r.meta.logger.Println("config file has new mod time, loading routes again...")

		err := r.file.Load(&newData)
		if err != nil {
			return fmt.Errorf("failed to get file modification time: %w", err)
		}

		r.Routes = newData.Routes

		err = r.file.StoreModTime()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Routes) LookupUrlFromPath(path string) (string, error) {
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

	var result any = r.Routes

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
		}

		if url, ok := result.(string); ok && strings.HasSuffix(url, "/*") {
			trailingKeys := keys[i+1:]

			url = strings.ReplaceAll(url, "/*", "")

			if len(trailingKeys) > 0 {
				result = strings.Join(append([]string{url}, trailingKeys...), "/")
			} else {
				result = url
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
