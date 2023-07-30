/*
Copyright 2022 Bill Nixon

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

// Config represents the configuration values.
type Config struct {
	Title               string // title of the application
	ServerHost          string // host to listen on
	ServerPort          string // port to listen on
	BaseURL             string // URL for password reset, e.g., https://host:port
	SQLDriverName       string // driverName for sql.Open
	SQLDataSourceName   string // dataSourceName for sql.Open
	ParseGlobPattern    string // pattern to use with template.ParseGlob
	SessionExpiresHours int    // number of hours session is valid
	SMTPHost            string // SMTP host to send email
	SMTPPort            string // SMTP port to send email
	SMTPUser            string // SMTP user to send email
	SMTPPassword        string // SMTP password to send email
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
	missing = appendIfEmpty(missing, c.ServerHost, "ServerHost")
	missing = appendIfEmpty(missing, c.ServerPort, "ServerPort")
	missing = appendIfEmpty(missing, c.BaseURL, "BaseURL")
	missing = appendIfEmpty(missing, c.SQLDriverName, "SQLDriverName")
	missing = appendIfEmpty(missing, c.SQLDataSourceName, "DataSourceName")
	missing = appendIfEmpty(missing, c.ParseGlobPattern, "ParseGlobPattern")
	missing = appendIfEmpty(missing, c.SMTPHost, "SMTPHost")
	missing = appendIfEmpty(missing, c.SMTPPort, "SMTPPort")
	missing = appendIfEmpty(missing, c.SMTPUser, "SMTPUser")
	missing = appendIfEmpty(missing, c.SMTPPassword, "SMTPPassword")

	return len(missing) == 0, missing
}

// RedactedConfig is a copy of Config used to redact values on output.
type RedactedConfig Config

// MarshalJSON is a custom Marshaler to redact some fields.
func (c Config) MarshalJSON() ([]byte, error) {
	r := RedactedConfig(c)
	r.SQLDataSourceName = "[REDACTED]"
	r.SMTPPassword = "[REDACTED]"
	return json.Marshal(r)
}

// String is a custom Stringer to redact some fields.
func (c Config) String() string {
	r := RedactedConfig(c)
	r.SQLDataSourceName = "[REDACTED]"
	r.SMTPPassword = "[REDACTED]"
	return fmt.Sprintf("%+v", r)
}
