package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"strings"
)

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

	// init logging
	err = InitLog(logFileName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w", err)
	}

	// read config file
	app.config, err = NewConfigFromFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w", err)
	}

	// ensure required config values have been provided
	isValid, missing := app.config.IsValid()
	if !isValid {
		return nil,
			fmt.Errorf("NewApp: invalid config: missing %s",
				strings.Join(missing, ", "))
	}

	// TODO: handle this default value
	if app.config.SessionExpiresHours == 0 {
		app.config.SessionExpiresHours = 24
	}

	// init database connection
	app.db, err = InitDB(app.config.SQLDriverName,
		app.config.SQLDataSourceName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w", err)
	}

	// init HTML templates
	app.tmpls, err = InitTemplates(app.config.ParseGlobPattern)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w", err)
	}

	return &app, err
}
