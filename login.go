/*
   Copyright 2022 Bill Nixon

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// LoginHandler handles /login requests.
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.tmpls, w, "login.html", nil)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}

	case http.MethodPost:
		app.loginPost(w, r)
	}
}

const (
	MsgMissingUserNameAndPassword = "Missing username and password"
	MsgMissingUserName            = "Missing username"
	MsgMissingPassword            = "Missing password"
	MsgLoginFailed                = "Login Failed"
)

// loginPost is called for the POST method of the LoginHandler.
func (app *App) loginPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))

	// check for missing values
	var msg string
	switch {
	case userName == "" && password == "":
		msg = MsgMissingUserNameAndPassword
	case userName == "":
		msg = MsgMissingUserName
	case password == "":
		msg = MsgMissingPassword
	}
	if msg != "" {
		log.Println(msg)
		err := RenderTemplate(app.tmpls, w, "login.html", msg)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// attempt to login the given userName with the given password
	token, err := app.LoginUser(userName, password)
	if err != nil {
		log.Printf("failed login for %q: %v", userName, err)
		err := RenderTemplate(app.tmpls, w, "login.html", MsgLoginFailed)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// login successful, so create a cookie for the session Token
	log.Printf("successful login for %q", userName)
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   token.Value,
		Expires: token.Expires,
	})

	// redirect from login page
	http.Redirect(w, r, "/hello", http.StatusSeeOther)
}

// LoginUser returns a session Token if userName and password is correct.
func (app *App) LoginUser(userName, password string) (Token, error) {
	err := CompareUserPassword(app.db, userName, password)
	if err != nil {
		return Token{}, errors.New("invalid password")
	}

	// create and save a new session token
	token, err := SaveNewToken(app.db, "session", userName, 32, app.config.SessionExpiresHours)
	if err != nil {
		log.Printf("unable to save session token: %v", err)
		return Token{}, fmt.Errorf("unable to save token: %w", err)
	}

	return token, nil
}
