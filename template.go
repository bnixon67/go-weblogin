package main

import (
	"html/template"
	"log"
)

func noescape(s string) template.HTML {
	return template.HTML(s)
}

// initTemplates parses the templates
func initTemplates(pattern string) (*template.Template, error) {
	log.Print("Initialize templates")

	return template.New("").Funcs(template.FuncMap{"noescape": noescape}).ParseGlob(pattern)
}
