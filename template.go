package main

import (
	"fmt"
	"html/template"
)

// InitTemplates parses the templates.
func InitTemplates(pattern string) (*template.Template, error) {
	tmpls, err := template.New("html").ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("InitTemplates: %w", err)
	}
	return tmpls, nil
}
