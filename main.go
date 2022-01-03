package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

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
		log.Printf("could not create NewApp, %v", err)
		return
	}

	// define HTTP server
	s := &http.Server{
		Addr:              ":" + app.config.ServerPort,
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
	http.HandleFunc("/w3.css",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "html/w3.css")
		})
	http.HandleFunc("/favicon.ico",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "html/favicon.ico")
		})
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
	var ip string
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}

	log.Println(ip, r.Method, r.RequestURI)
	l.next.ServeHTTP(w, r)
}
