package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

	if len(os.Args) != 2 {
		fmt.Printf("%s [CONFIG FILE]\n", os.Args[0])
		return
	}

	// read config file
	configFileName := os.Args[1]
	if configFileName == "" {
		fmt.Printf("%s [CONFIG FILE]\n", os.Args[0])
		return
	}
	config, err = readConfig(configFileName)
	if err != nil {
		log.Printf("Unable to read config file %q", configFileName)
		log.Panic(err)
	}

	// ensure required config values have been provided
	logPanicIfEmpty(config.SQLDriverName, "Missing SQLDriverName in config file")
	logPanicIfEmpty(config.SQLDataSourceName, "Missing SQLDataSourceName in config file")
	logPanicIfEmpty(config.ParseGlobPattern, "Missing ParseGlobPattern in config file")
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
	s := &http.Server{
		Addr:              ":8443",
		Handler:           &logRequestHandler{http.DefaultServeMux},
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// register handlers
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/forgot", ForgotHandler)
	http.HandleFunc("/reset", ResetHandler)
	http.HandleFunc("/hello", HelloHandler)
	http.Handle("/style.css", http.FileServer(http.Dir("html")))
	http.Handle("/", http.RedirectHandler("/hello", http.StatusMovedPermanently))

	// run server
	// TODO: move certs to config file
	log.Panic(s.ListenAndServeTLS("cert/cert.pem", "cert/key.pem"))
}

type logRequestHandler struct {
	next http.Handler
}

func (l *logRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI, r.Header)
	l.next.ServeHTTP(w, r)
}
