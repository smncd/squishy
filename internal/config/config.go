package config

import (
	"os"

	"github.com/alexflint/go-arg"
)

type Options struct {
	Debug bool   `yaml:"debug" json:"debug" arg:"-D,--debug,env:DEBUG"`
	Host  string `yaml:"host" json:"host" validate:"required" arg:"-H,--host,env:HOST"`
	Port  int    `yaml:"port" json:"port" validate:"required" arg:"-P,--port,env:PORT"`
}

type Config struct {
	file
	Options `yaml:"config" json:"config"`
}

func New() (*Config, error) {
	args := os.Args[1:]

	config := Config{
		file: file{
			Path: "squishy.yaml",
		},
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

	err = config.file.Load(&config)
	if err != nil {
		return nil, err
	}

	// we parse a second time to overwrite
	// file-based config with flags (if applicable)
	parser.MustParse(args)

	return &config, nil
}
