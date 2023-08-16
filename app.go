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

var (
	ErrAppGetConfig     = errors.New("failed")
	ErrAppInvalidConfig = errors.New("invalid config")
	ErrAppInitDB        = errors.New("failed")
	ErrAppInitTemplates = errors.New("failed")
)

// App contains common variables to avoid using global variables.
type App struct {
	DB    *sql.DB
	Tmpls *template.Template
	Cfg   Config
}

// NewApp returns a new App based on the config filename provided.
func NewApp(configFilename string) (*App, error) {
	fn := "NewApp"

	var app App
	var err error

	// read config file
	app.Cfg, err = GetConfigFromFile(configFilename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", fn, ErrAppGetConfig, err)
	}

	// ensure required config values have been provided
	isValid, missing := app.Cfg.IsValid()
	if !isValid {
		return nil, fmt.Errorf("%s: %w: missing %s", fn, ErrAppInvalidConfig, strings.Join(missing, ", "))
	}

	// default to 24 hours if no session expiration
	if app.Cfg.SessionExpiresHours == 0 {
		app.Cfg.SessionExpiresHours = 24
	}

	// init database connection
	app.DB, err = InitDB(app.Cfg.SQL.DriverName, app.Cfg.SQL.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", fn, ErrAppInitDB, err)
	}

	// init HTML templates
	app.Tmpls, err = InitTemplates(app.Cfg.ParseGlobPattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", fn, ErrAppInitTemplates, err)
	}

	return &app, err
}
