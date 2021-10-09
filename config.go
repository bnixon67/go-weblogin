package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config represents the configuration values
type Config struct {
	// SQLDriverName is the driverName to use with sql.Open
	SQLDriverName string

	// SQLDataSourceName is the dataSourceName to use with sql.Open
	SQLDataSourceName string

	// ParseGlobPattern is the pattern to use with template.ParseGlob
	ParseGlobPattern string

	// SessionExpiresHours is the number of hours after a session expires
	SessionExpiresHours int

	SmtpHost     string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
}

// NewConfigFromFile returns a Config from the given fileName
func NewConfigFromFile(fileName string) (Config, error) {
	log.Printf("reading %q", fileName)

	config := Config{}

	configFile, err := os.Open(fileName)
	if err != nil {
		return config, err
	}
	defer closeConfig(configFile)

	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}

// closeConfig closes the Config file
func closeConfig(f *os.File) {
	log.Printf("closing %q", f.Name())

	err := f.Close()
	if err != nil {
		log.Panic(err)
	}
}

// IsValid returns true if the config is valid
func (c Config) IsValid() bool {
	// ensure required config values have been provided
	// test conditions on separate lines to avoid short circuit evaluation
	isEmpty := false
	isEmpty = logIfEmpty(c.SQLDriverName, "missing or empty SQLDriverName in config file") || isEmpty
	isEmpty = logIfEmpty(c.SQLDataSourceName, "missing or empty SQLDataSourceName in config file") || isEmpty
	isEmpty = logIfEmpty(c.ParseGlobPattern, "missing or empty ParseGlobPattern in config file") || isEmpty

	return !isEmpty
}
