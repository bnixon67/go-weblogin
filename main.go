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

type App struct {
	db     *sql.DB
	tmpls  *template.Template
	config Config
}

func NewApp(configFileName, logFileName string) (*App, error) {
	var app App
	var err error

	// use custom writer for log
	lw, err := NewLogWriter(logFileName)
	if err != nil {
		log.Printf("unable to create NewLogWriter, %v", err)
		return &app, err
	}
	log.SetFlags(0)
	log.SetOutput(lw)

	// read config file
	app.config, err = NewConfigFromFile(configFileName)
	if err != nil {
		log.Printf("unable to read config file %q, %v", configFileName, err)
		return &app, err
	}

	// ensure required config values have been provided
	if !app.config.IsValid() {
		log.Printf("config is not valid")
		return &app, err
	}

	// TODO: handle this default value
	if app.config.SessionExpiresHours == 0 {
		app.config.SessionExpiresHours = 24
	}

	// init database connection
	app.db, err = initDB(app.config.SQLDriverName, app.config.SQLDataSourceName)
	if err != nil {
		log.Printf("initDB failed: %v", err)
		return &app, err
	}

	// init HTML templates
	app.tmpls, err = initTemplates(app.config.ParseGlobPattern)
	if err != nil {
		log.Printf("initTemplates failed: %v", err)
		return &app, err
	}

	return &app, err
}

// main function
func main() {
	// config file must be passed as argument and not empty
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Printf("%s [CONFIG FILE]\n", os.Args[0])
		return
	}

	// TODO: allow logfile to specified in config file
	app, err := NewApp(os.Args[1], "")
	if err != nil {
		log.Println("init failed", err)
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
	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/register", app.RegisterHandler)
	http.HandleFunc("/logout", app.LogoutHandler)
	http.HandleFunc("/forgot", app.ForgotHandler)
	http.HandleFunc("/reset", app.ResetHandler)
	http.HandleFunc("/hello", app.HelloHandler)
	http.Handle("/style.css", http.FileServer(http.Dir("html")))
	http.Handle("/", http.RedirectHandler("/hello", http.StatusMovedPermanently))

	// run server
	// TODO: move certs to config file
	log.Println("Listening on", s.Addr)
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
