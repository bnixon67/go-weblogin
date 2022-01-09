package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoSuchUser      = errors.New("no such user")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInternalFailure = errors.New("login failed due to internal error")
)

// LoginPageData record
type LoginPageData struct {
	Message     string
	CurrentUser User
}

// LoginHandler handles /login requests
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// get currentUser if sessionToken exists
		var sessionToken string
		var currentUser User
		var err error
		c, err := r.Cookie("sessionToken")
		if err == nil {
			sessionToken = c.Value
			currentUser, _ = app.GetUserForSessionToken(sessionToken)
			log.Printf("%+v", currentUser)
		}

		err = app.tmpls.ExecuteTemplate(w, "login.html", nil)
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
	userName := r.PostFormValue("username")
	password := r.PostFormValue("password")

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
	sessionToken, sessionExpires, err := app.loginUser(userName, password)
	if err != nil {
		err := app.tmpls.ExecuteTemplate(w, "login.html", MSG_LOGIN_FAILED)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// login successful

	// create a cookie for the sessionToken
	http.SetCookie(w, &http.Cookie{
		Name:    "sessionToken",
		Value:   sessionToken,
		Expires: sessionExpires,
	})
	log.Printf("Login for %q successful with sessionToken %q expires %q",
		userName, sessionToken, sessionExpires.UTC().Format(time.RFC3339Nano))

	http.Redirect(w, r, "/hello", http.StatusSeeOther)
}

// loginUser returns a sessionToken if the given userName and password is correct, otherwise error
func (app *App) loginUser(userName, password string) (string, time.Time, error) {
	var sessionToken string
	var sessionExpires time.Time

	// get hashed password for the given user
	result := app.db.QueryRow("SELECT hashedPassword FROM users WHERE username=?", userName)
	var hashedPassword string
	err := result.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("%v %q", ErrNoSuchUser, userName)
			return sessionToken, sessionExpires, ErrNoSuchUser
		}
		log.Printf("query failed for %q: %v", userName, err)
		return sessionToken, sessionExpires, ErrInternalFailure
	}

	// compared hashed password with given password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("%v for %q", ErrInvalidPassword, userName)
		return sessionToken, sessionExpires, ErrInvalidPassword
	}

	// create a new random sessions token
	sessionToken, err = GenerateRandomString(32)
	if err != nil {
		log.Printf("unable to GenerateRandomString: %v", err)
		return sessionToken, sessionExpires, ErrInternalFailure
	}

	// store the sessionToken
	sessionExpires = time.Now().Add(time.Duration(app.config.SessionExpiresHours) * time.Hour)

	_, err = app.db.Query("UPDATE users SET sessionToken = ?, sessionExpires = ? WHERE username = ?", sessionToken, sessionExpires, userName)
	if err != nil {
		log.Printf("update sessionToken failed for %q: %v", userName, err)
		return sessionToken, sessionExpires, ErrInternalFailure
	}

	return sessionToken, sessionExpires, nil
}
