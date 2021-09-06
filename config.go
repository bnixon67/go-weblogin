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
}

// readConfig return the Config from the given fileName
func readConfig(fileName string) (Config, error) {
	log.Print("reading config file")

	configFile, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer closeConfig(configFile)

	config := Config{}
	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}

// closeConfig closes the Config file
func closeConfig(f *os.File) {
	log.Println("closing config file")

	err := f.Close()
	if err != nil {
		log.Panic(err)
	}
}
