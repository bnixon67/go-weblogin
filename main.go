package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// db is the global datbase handle
var db *sql.DB

// tmpls is the gloabl for the parsed HTML templates
var tmpls *template.Template

// config is the global for the config values
var config Config

// main function
func main() {
	var err error

	// use custom log writer
	log.SetFlags(0)
	log.SetOutput(new(LogWriter))

	// read config file
	configFileName := "config.json"
	config, err = readConfig(configFileName)
	if err != nil {
		log.Printf("Unable to read config file %q", configFileName)
		log.Panic(err)
	}

	// ensure required config values have been provided
	logPanicIsEmpty(config.SQLDriverName, "Missing SQLDriverName in config file")
	logPanicIsEmpty(config.SQLDataSourceName, "Missing SQLDataSourceName in config file")
	logPanicIsEmpty(config.ParseGlobPattern, "Missing ParseGlobPattern in config file")
	if config.SessionExpiresHours == 0 {
		config.SessionExpiresHours = 24
	}

	// init database connection
	db, err = initDB(config.SQLDriverName, config.SQLDataSourceName)
	if err != nil {
		log.Panic(err)
	}

	// init HTML templates
	tmpls, err = initTemplates(config.ParseGlobPattern)
	if err != nil {
		log.Panic(err)
	}

	// define HTTP server
	// TODO: move values to config file
	s := &http.Server{
		Addr:           ":8000",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// register handlers
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/hello", HelloHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.Handle("/style.css", http.FileServer(http.Dir("html")))

	// run server
	log.Panic(s.ListenAndServe())
}
