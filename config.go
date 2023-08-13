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
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ConfigSQL struct {
	DriverName     string // driverName for sql.Open
	DataSourceName string // dataSourceName for sql.Open
}

type ConfigSMTP struct {
	Host     string // SMTP host to send email
	Port     string // SMTP port to send email
	User     string // SMTP user to send email
	Password string // SMTP password to send email
}

type ConfigServer struct {
	Host string // host to listen on
	Port string // port to listen on
}

// Config represents the configuration values.
type Config struct {
	Title               string // title of the application
	BaseURL             string // base URL, e.g., https://host:port
	ParseGlobPattern    string // pattern to use with template.ParseGlob
	SessionExpiresHours int    // number of hours session is valid
	Server              ConfigServer
	SQL                 ConfigSQL
	SMTP                ConfigSMTP
}

var (
	ErrConfigOpen   = errors.New("failed to open")
	ErrConfigDecode = errors.New("failed to decode")
)

// NewConfigFromFile returns a Config from the given fileName.
func NewConfigFromFile(fileName string) (Config, error) {
	var config Config

	// open config file
	configFile, err := os.Open(fileName)
	if err != nil {
		return config, fmt.Errorf("NewConfigFromFile: %w: %v", ErrConfigOpen, err)
	}
	defer configFile.Close()

	// decode json from config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("NewConfigFromFile: %w: %v", ErrConfigDecode, err)
	}

	return config, nil
}

// appendIfEmpty appends msg to target if str is empty and returns target.
func appendIfEmpty(target []string, str, msg string) []string {
	if str == "" {
		target = append(target, msg)
	}

	return target
}

// IsValid returns true if the config has all the required values.
func (c *Config) IsValid() (bool, []string) {
	var missing []string

	missing = appendIfEmpty(missing, c.Title, "Title")
	missing = appendIfEmpty(missing, c.BaseURL, "BaseURL")
	missing = appendIfEmpty(missing, c.ParseGlobPattern, "ParseGlobPattern")
	missing = appendIfEmpty(missing, c.Server.Host, "Server.Host")
	missing = appendIfEmpty(missing, c.Server.Port, "Server.Port")
	missing = appendIfEmpty(missing, c.SQL.DriverName, "SQL.DriverName")
	missing = appendIfEmpty(missing, c.SQL.DataSourceName, "SQL.DataSourceName")
	missing = appendIfEmpty(missing, c.SMTP.Host, "SMTP.Host")
	missing = appendIfEmpty(missing, c.SMTP.Port, "SMTP.Port")
	missing = appendIfEmpty(missing, c.SMTP.User, "SMTP.User")
	missing = appendIfEmpty(missing, c.SMTP.Password, "SMTP.Password")

	return len(missing) == 0, missing
}

// RedactedConfig is a copy of Config used to redact values on output.
type RedactedConfig Config

// MarshalJSON is a custom Marshaler to redact some fields.
func (c Config) MarshalJSON() ([]byte, error) {
	r := RedactedConfig(c)
	r.SQL.DataSourceName = "[REDACTED]"
	r.SMTP.Password = "[REDACTED]"
	return json.Marshal(r)
}

// String is a custom Stringer to redact some fields.
func (c Config) String() string {
	r := RedactedConfig(c)
	r.SQL.DataSourceName = "[REDACTED]"
	r.SMTP.Password = "[REDACTED]"
	return fmt.Sprintf("%+v", r)
}
