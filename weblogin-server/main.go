/*
Copyright 2023 Bill Nixon

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
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	weblogin "github.com/bnixon67/go-weblogin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// define command-line flags
	configFilename := flag.String("config", "", "config filename")
	logFilename := flag.String("log", "", "log filename")
	logLevel := flag.Int("logLevel", 0, "log level")
	logAddSource := flag.Bool("logAddSource", false, "log add sourc")

	// define custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// parse command-line flags
	flag.Parse()

	// configFilename is required
	if *configFilename == "" {
		flag.Usage()
		os.Exit(2)
	}

	// check for additional command-line arguments
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(3)
	}

	weblogin.InitLog(*logFilename, slog.Level(*logLevel), *logAddSource)

	app, err := weblogin.NewApp(*configFilename)
	if err != nil {
		slog.Error("failed to create app", "err", err)
		return
	}
	slog.Info("created app", "app", app)

	mux := http.NewServeMux()

	// define HTTP server
	// TODO: add values to config file
	srv := &http.Server{
		Addr:              ":" + app.Config.Server.Port,
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
		slog.Info("starting server",
			slog.Group("srv",
				"Addr", srv.Addr,
				"ReadTimeout (s)", srv.ReadTimeout/time.Second,
				"WriteTimeout (s)", srv.WriteTimeout/time.Second,
				"IdleTimeout (s)", srv.IdleTimeout/time.Second,
				"ReadHeaderTimeout (s)", srv.ReadHeaderTimeout/time.Second,
				"MaxHeaderBytes (kb)", srv.MaxHeaderBytes/1024,
			),
		)
		// TODO: move cert locations to config file
		err = srv.ListenAndServeTLS("cert/cert.pem", "cert/key.pem")
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "err", err)
			os.Exit(1)
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
		slog.Error("server shutdown error", "err", err)
	}

	slog.Info("server closed")
}
