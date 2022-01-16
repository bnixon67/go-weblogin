package main

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"strings"
)

var ErrInvalidConfig = errors.New("invalid config")

// App contains common variables to reuse to eliminate global variables.
type App struct {
	db     *sql.DB
	tmpls  *template.Template
	config Config
}

// NewApp returns a new App based on the config and log filenames provided.
func NewApp(configFileName, logFileName string) (*App, error) {
	var app App
	var err error

	err = InitLogging(logFileName)
	if err != nil {
		return nil, err
	}

	// read config file
	app.config, err = NewConfigFromFile(configFileName)
	if err != nil {
		log.Printf("unable to read config file %q, %v", configFileName, err)
		return nil, err
	}

	// ensure required config values have been provided
	isValid, missing := app.config.IsValid()
	if !isValid {
		log.Printf("config missing %s", strings.Join(missing, ", "))
		return nil, ErrInvalidConfig
	}

	// TODO: handle this default value
	if app.config.SessionExpiresHours == 0 {
		app.config.SessionExpiresHours = 24
	}

	// init database connection
	app.db, err = initDB(app.config.SQLDriverName, app.config.SQLDataSourceName)
	if err != nil {
		log.Printf("initDB failed: %v", err)
		return nil, err
	}

	// init HTML templates
	app.tmpls, err = initTemplates(app.config.ParseGlobPattern)
	if err != nil {
		log.Printf("initTemplates failed: %v", err)
		return nil, err
	}

	return &app, err
}
