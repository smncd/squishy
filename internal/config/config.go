package config

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Debug bool   `yaml:"debug" json:"debug" arg:"-D,--debug,env:DEBUG"`
	Host  string `yaml:"host" json:"host" validate:"required" arg:"-H,--host,env:HOST"`
	Port  int    `yaml:"port" json:"port" validate:"required" arg:"-P,--port,env:PORT"`
}

type Config struct {
	Options `yaml:"config" json:"config"`
}

func New() (*Config, error) {
	args := os.Args[1:]

	config := Config{
		Options: Options{
			Debug: false,
			Host:  "localhost",
			Port:  1394,
		},
	}

	parser, err := arg.NewParser(arg.Config{EnvPrefix: "SQUISHY_"}, &config)
	if err != nil {
		return nil, err
	}

	parser.MustParse(args)

	err = config.loadConfigFromFile()
	if err != nil {
		return nil, err
	}

	// we parse a second time to overwrite
	// file-based config with flags (if applicable)
	parser.MustParse(args)

	return &config, nil
}

func (c *Config) loadConfigFromFile() error {
	filePath := "squishy.yaml"

	rawData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file (%s): %w", filePath, err)
	}

	if err := yaml.Unmarshal(rawData, &c); err != nil {
		return fmt.Errorf("failed to unmarshal YAML config: %w", err)
	}

	return nil
}
