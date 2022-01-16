package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config represents the configuration values.
type Config struct {
	ServerHost string // host to listen on
	ServerPort string // port to listen on

	SQLDriverName     string // driverName for sql.Open
	SQLDataSourceName string // dataSourceName for sql.Open

	ParseGlobPattern string // pattern to use with template.ParseGlob

	SessionExpiresHours int // number of hours session is valid

	SmtpHost     string // SMTP host to send email
	SmtpPort     string // SMTP port to send email
	SmtpUser     string // SMTP user to send email
	SmtpPassword string // SMTP password to send email
}

// NewConfigFromFile returns a Config from the given fileName.
func NewConfigFromFile(fileName string) (Config, error) {
	log.Printf("INFO - reading %q", fileName)

	configFile, err := os.Open(fileName)
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()

	var config Config
	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}

func appendIfEmpty(missing []string, str, msg string) []string {
	if str == "" {
		missing = append(missing, msg)
	}

	return missing
}

// IsValid returns true if the config has all the required values.
func (c Config) IsValid() (bool, []string) {
	var missing []string

	missing = appendIfEmpty(missing, c.ServerHost, "ServerHost")
	missing = appendIfEmpty(missing, c.ServerPort, "ServerPort")
	missing = appendIfEmpty(missing, c.SQLDriverName, "SQLDriverName")
	missing = appendIfEmpty(missing, c.SQLDataSourceName, "DataSourceName")
	missing = appendIfEmpty(missing, c.ParseGlobPattern, "ParseGlobPattern")

	return missing == nil, missing
}
