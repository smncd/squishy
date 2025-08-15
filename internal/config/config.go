package config

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"gitlab.com/smncd/squishy/internal/logging"
)

type Options struct {
	Debug bool   `yaml:"debug" json:"debug" arg:"-D,--debug,env:DEBUG"`
	Host  string `yaml:"host" json:"host" validate:"required" arg:"-H,--host,env:HOST"`
	Port  int    `yaml:"port" json:"port" validate:"required" arg:"-P,--port,env:PORT"`
}

type Config struct {
	file
	Options `yaml:"config" json:"config"`
	logger  *log.Logger
}

func New(logger *log.Logger) (*Config, error) {
	args := os.Args[1:]

	config := Config{
		Options: Options{
			Debug: false,
			Host:  "localhost",
			Port:  1394,
		},
		logger: logger,
	}

	parser, err := arg.NewParser(arg.Config{EnvPrefix: "SQUISHY_"}, &config)
	if err != nil {
		return nil, err
	}

	parser.MustParse(args)

	logging.Info(logger, "Loading config file with path: %s", config.file.Path)
	err = config.file.Load(&config)
	if err != nil {
		return nil, err
	}

	// we parse a second time to overwrite
	// file-based config with flags (if applicable)
	parser.MustParse(args)

	return &config, nil
}

func (cfg *Config) Routes() (*Routes, error) {
	routes := &Routes{
		file:   cfg.file,
		logger: cfg.logger,
	}

	err := routes.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading routes: %w", err)
	}

	return routes, nil
}
