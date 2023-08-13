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

var (
	ErrConfigOpen   = errors.New("failed")
	ErrConfigDecode = errors.New("failed to decode")
)

// ConfigSQL contains SQL related configuration values.
type ConfigSQL struct {
	DriverName     string
	DataSourceName string
}

// ConfigSMTP contains SMTP related configuration values.
type ConfigSMTP struct {
	Host     string
	Port     string
	User     string
	Password string
}

// ConfigServer contains Server related configuration values.
type ConfigServer struct {
	Host string
	Port string
}

// Config represents the configuration values.
type Config struct {
	Title               string // title of the application
	BaseURL             string // base URL, e.g., https://host:port
	ParseGlobPattern    string // pattern to use with template.ParseGlob
	SessionExpiresHours int    // number of hours user session is valid
	Server              ConfigServer
	SQL                 ConfigSQL
	SMTP                ConfigSMTP
}

// GetConfigFromFile returns the Config from filename.
func GetConfigFromFile(filename string) (Config, error) {
	fn := "GetConfigFromFile"

	var config Config

	// open config file
	configFile, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("%s: %w: %v", fn, ErrConfigOpen, err)
	}
	defer configFile.Close()

	// decode json from config
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, fmt.Errorf("%s: %w: %v", fn, ErrConfigDecode, err)
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
