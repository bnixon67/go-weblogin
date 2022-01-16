package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInternalFailure = errors.New("login failed due to internal error")
)

// LoginHandler handles /login requests
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := app.tmpls.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodPost:
		app.loginPost(w, r)
	}
}

const (
	MSG_MISSING_USERNAME_PASSWORD = "Missing username and password"
	MSG_MISSING_USERNAME          = "Missing username"
	MSG_MISSING_PASSWORD          = "Missing password"
	MSG_LOGIN_FAILED              = "Login Failed"
)

// loginPost is called for the POST method of the LoginHandler
func (app *App) loginPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))

	// check for missing values
	var msg string
	switch {
	case userName == "" && password == "":
		msg = MSG_MISSING_USERNAME_PASSWORD
	case userName == "":
		msg = MSG_MISSING_USERNAME
	case password == "":
		msg = MSG_MISSING_PASSWORD
	}
	if msg != "" {
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, "login.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// attempt to login the given userName with the given password
	token, err := app.LoginUser(userName, password)
	if err != nil {
		err := app.tmpls.ExecuteTemplate(w, "login.html", MSG_LOGIN_FAILED)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// login successful, so create a cookie for the session Token
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   token.Value,
		Expires: token.Expires,
	})
	log.Printf("valid login for %q", userName)

	http.Redirect(w, r, "/hello", http.StatusSeeOther)
}

// LoginUser returns a session Token if userName and password is correct
func (app *App) LoginUser(userName, password string) (Token, error) {
	err := app.CheckUserPassword(userName, password)
	if err != nil {
		log.Printf("invalid password for %q: %v", userName, err)
		return Token{}, err
	}

	// create and save a new session token
	token, err := SaveNewToken(app.db, "session", userName, app.config.SessionExpiresHours)
	if err != nil {
		log.Printf("unable to save session token: %v", err)
		return Token{}, ErrInternalFailure
	}

	return token, nil
}
