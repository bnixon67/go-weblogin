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

	// config file must be passed as argument and not empty
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Printf("%s [CONFIG FILE]\n", os.Args[0])
		return
	}
	configFileName := os.Args[1]

	// read config file
	config, err = NewConfigFromFile(configFileName)
	if err != nil {
		log.Printf("unable to read config file %q, %v", configFileName, err)
		return
	}

	// ensure required config values have been provided
	if !config.IsValid() {
		log.Printf("config is not valid")
		return
	}

	// TODO: handle this default value
	if config.SessionExpiresHours == 0 {
		config.SessionExpiresHours = 24
	}

	// init database connection
	db, err = initDB(config.SQLDriverName, config.SQLDataSourceName)
	if err != nil {
		log.Printf("initDB failed: %v", err)
		return
	}

	// init HTML templates
	tmpls, err = initTemplates(config.ParseGlobPattern)
	if err != nil {
		log.Printf("initTemplates failed: %v", err)
		return
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
	err = s.ListenAndServeTLS("cert/cert.pem", "cert/key.pem")
	if err != nil {
		log.Printf("ListandServeTLS failed: %v", err)
	}
}

// logRequestHandler is middleware that logs all HTTP requests and then calls the next HTTP handler specified
type logRequestHandler struct {
	next http.Handler
}

// ServerHTTP for logRequestHandler log the HTTP request and then calls the next HTTP handler specified
func (l *logRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.Method, r.RequestURI)
	l.next.ServeHTTP(w, r)
}
