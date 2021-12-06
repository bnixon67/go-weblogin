package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config represents the configuration values.
type Config struct {
	ServerAddr string // server address

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

// closeConfig closes the Config file. This is mostly unnecessary for readable files, but may help in debugging. See https://www.joeshaw.org/dont-defer-close-on-writable-files/ for more information.
func closeConfig(f *os.File) {
	if f == nil {
		log.Printf("unexpected nil file")
		return
	}

	log.Printf("closing %q", f.Name())

	err := f.Close()
	if err != nil {
		log.Print("Close() failed, ", err)
	}
}

// IsValid returns true if the config has all the required values.
func (c Config) IsValid() bool {
	// ensure required config values have been provided
	// test conditions on separate lines to avoid short circuit evaluation
	isEmpty := false
	isEmpty = logIfEmpty(c.ServerAddr, "missing or empty ServerAddr in config file") || isEmpty
	isEmpty = logIfEmpty(c.SQLDriverName, "missing or empty SQLDriverName in config file") || isEmpty
	isEmpty = logIfEmpty(c.SQLDataSourceName, "missing or empty SQLDataSourceName in config file") || isEmpty
	isEmpty = logIfEmpty(c.ParseGlobPattern, "missing or empty ParseGlobPattern in config file") || isEmpty

	return !isEmpty
}
