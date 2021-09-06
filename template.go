package main

import (
	"html/template"
	"log"
)

// initTemplates parses the templates
func initTemplates(pattern string) (*template.Template, error) {
	log.Print("Initialize templates")

	return template.ParseGlob(pattern)
}
