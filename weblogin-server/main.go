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
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	weblogin "github.com/bnixon67/go-web-login"
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

	// define HTTP server
	// TODO: add values to config file
	s := &http.Server{
		Addr: ":" + app.Config.ServerPort,
		Handler: &weblogin.LogRequestHandler{
			Next: http.DefaultServeMux,
		},
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
	// TODO: define base html directory in config
	http.HandleFunc("/w3.css", weblogin.ServeFileHandler("../html/w3.css"))
	http.HandleFunc("/favicon.ico", weblogin.ServeFileHandler("../html/favicon.ico"))
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
