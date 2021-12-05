package main

import (
	"database/sql"
	"html/template"
	"log"
)

type App struct {
	db     *sql.DB
	tmpls  *template.Template
	config Config
}

func NewApp(configFileName, logFileName string) (*App, error) {
	var app App
	var err error

	err = InitLogging(logFileName)
	if err != nil {
		return &app, err
	}

	// read config file
	app.config, err = NewConfigFromFile(configFileName)
	if err != nil {
		log.Printf("unable to read config file %q, %v", configFileName, err)
		return &app, err
	}

	// ensure required config values have been provided
	if !app.config.IsValid() {
		log.Printf("config is not valid")
		return &app, err
	}

	// TODO: handle this default value
	if app.config.SessionExpiresHours == 0 {
		app.config.SessionExpiresHours = 24
	}

	// init database connection
	app.db, err = initDB(app.config.SQLDriverName, app.config.SQLDataSourceName)
	if err != nil {
		log.Printf("initDB failed: %v", err)
		return &app, err
	}

	// init HTML templates
	app.tmpls, err = initTemplates(app.config.ParseGlobPattern)
	if err != nil {
		log.Printf("initTemplates failed: %v", err)
		return &app, err
	}

	return &app, err
}
