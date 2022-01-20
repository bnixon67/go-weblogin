package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// config file must be passed as argument and not empty
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Printf("%s [CONFIG FILE]\n", os.Args[0])
		return
	}

	// TODO: allow logfile to specified in config file
	configFileName := os.Args[1]
	logFileName := ""
	app, err := NewApp(configFileName, logFileName)
	if err != nil {
		log.Printf("failed to create app: %v", err)
		return
	}
	log.Printf("created app using config %q and log %q",
		configFileName, logFileName)

	// define HTTP server
	// TODO: add values to config file
	s := &http.Server{
		Addr:              ":" + app.config.ServerPort,
		Handler:           &LogRequestHandler{http.DefaultServeMux},
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
	http.HandleFunc("/w3.css", ServeFileHandler("html/w3.css"))
	http.HandleFunc("/favicon.ico", ServeFileHandler("html/favicon.ico"))
	http.Handle("/",
		http.RedirectHandler("/hello", http.StatusMovedPermanently))

	// run server
	// TODO: move cert locations to config file
	log.Println("Listening on", s.Addr)
	err = s.ListenAndServeTLS("cert/cert.pem", "cert/key.pem")
	if err != nil {
		log.Printf("ListandServeTLS failed: %v", err)
	}
}

// ServeFileHandler is a simple http.ServeFile wrapper.
func ServeFileHandler(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}

// LogRequestHandler is middleware that logs all HTTP requests and
// then calls the next HTTP handler specified.
type LogRequestHandler struct {
	next http.Handler
}

// ServerHTTP for logRequestHandler log the HTTP request and then
// calls the next HTTP handler specified.
func (l *LogRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get real IP address if using Cloudflare or similar service
	var ip string
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}

	log.Println(ip, r.Method, r.RequestURI)

	l.next.ServeHTTP(w, r)
}
