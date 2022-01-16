package main

import (
	"html/template"
	"log"
)

func noescape(s string) template.HTML {
	return template.HTML(s) //nolint
}

// InitTemplates parses the templates.
func InitTemplates(pattern string) (*template.Template, error) {
	log.Print("Initialize templates")

	return template.New("").Funcs(template.FuncMap{"noescape": noescape}).ParseGlob(pattern)
}
