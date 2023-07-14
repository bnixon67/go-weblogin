/*
Copyright 2022 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	weblogin "github.com/bnixon67/go-weblogin"
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
	app, err := weblogin.NewApp(configFileName, logFileName)
	if err != nil {
		log.Printf("failed to create app: %v", err)
		return
	}
	log.Printf("created app using config %q and log %q",
		configFileName, logFileName)

	mux := http.NewServeMux()

	// define HTTP server
	// TODO: add values to config file
	srv := &http.Server{
		Addr:              ":" + app.Config.ServerPort,
		Handler:           &weblogin.LogRequestHandler{Next: mux},
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	// register handlers
	mux.HandleFunc("/login", app.LoginHandler)
	mux.HandleFunc("/register", app.RegisterHandler)
	mux.HandleFunc("/logout", app.LogoutHandler)
	mux.HandleFunc("/forgot", app.ForgotHandler)
	mux.HandleFunc("/reset", app.ResetHandler)
	mux.HandleFunc("/hello", app.HelloHandler)
	mux.HandleFunc("/users", app.UsersHandler)
	// TODO: define base html directory in config
	mux.HandleFunc("/w3.css", weblogin.ServeFileHandler("../html/w3.css"))
	mux.HandleFunc("/favicon.ico", weblogin.ServeFileHandler("../html/favicon.ico"))
	mux.Handle("/",
		http.RedirectHandler("/hello", http.StatusMovedPermanently))

	// create a channel to receive interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// start the server in a goroutine
	go func() {
		log.Println("Listening on", srv.Addr)
		// TODO: move cert locations to config file
		err = srv.ListenAndServeTLS("cert/cert.pem", "cert/key.pem")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ListandServeTLS failed: %v", err)
		}
	}()

	// wait for an interrupt signal
	<-interrupt

	// create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// initiate the shutdown process
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Println("server shutdown error", "err", err)
	}

	log.Println("server closed")
}
