package resources

import "embed"

// Static assets and templates.

//go:embed templates/*
var TemplateFS embed.FS

//go:embed static/*
var StaticFS embed.FS

type ErrorTemplateData struct {
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Error       string
}
