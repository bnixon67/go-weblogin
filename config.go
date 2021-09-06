package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config defines the configuration values
type Config struct {
	// SqlDriverName is the driverName to use with sql.Open
	SqlDriverName string

	// SqlDataSourceName is the dataSourceName to use with sql.Open
	SqlDataSourceName string

	// ParseGlobPattern is the pattern to use with template.ParseGlob
	ParseGlobPattern string
}

// readConfig return the Config from the given fileName
func readConfig(fileName string) (Config, error) {
	log.Print("Reading config file")

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
	log.Println("Closing config file")

	err := f.Close()
	if err != nil {
		log.Panic(err)
	}
}
