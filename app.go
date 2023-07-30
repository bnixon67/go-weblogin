/*
Copyright 2023 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"strings"
)

// App contains common variables to avoid using global variables.
type App struct {
	DB     *sql.DB
	Tmpls  *template.Template
	Config Config
}

var (
	ErrInitLog       = errors.New("failed to init log")
	ErrOpenConfig    = errors.New("failed to open config")
	ErrInvalidConfig = errors.New("invalid config")
	ErrInitDB        = errors.New("failed to init db")
	ErrInitTemplates = errors.New("failed to init templates")
)

// NewApp returns a new App based on the config and log filenames provided.
func NewApp(configFileName, logFileName string) (*App, error) {
	var app App
	var err error

	// init log
	err = InitLog(logFileName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w: %v", ErrInitLog, err)
	}

	// read config file
	app.Config, err = NewConfigFromFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w: %v", ErrOpenConfig, err)
	}

	// ensure required config values have been provided
	isValid, missing := app.Config.IsValid()
	if !isValid {
		return nil, fmt.Errorf("NewApp: %w: missing %s", ErrInvalidConfig, strings.Join(missing, ", "))
	}

	// TODO: handle this default value
	if app.Config.SessionExpiresHours == 0 {
		app.Config.SessionExpiresHours = 24
	}

	// init database connection
	app.DB, err = InitDB(app.Config.SQLDriverName, app.Config.SQLDataSourceName)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w: %v", ErrInitDB, err)
	}

	// init HTML templates
	app.Tmpls, err = InitTemplates(app.Config.ParseGlobPattern)
	if err != nil {
		return nil, fmt.Errorf("NewApp: %w: %v", ErrInitTemplates, err)
	}

	return &app, err
}
