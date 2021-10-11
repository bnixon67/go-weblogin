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

// LoginPageData record
type LoginPageData struct {
	Message     string
	CurrentUser User
}

// LoginHandler handles /login requests
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.Method)

	switch r.Method {

	case "GET":
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

	case "POST":
		app.loginPut(w, r)

	default:
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// loginPut is called for the PUT method of the LoginHandler
func (app *App) loginPut(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := r.PostFormValue("username")
	password := r.PostFormValue("password")

	// check for missing values
	var msg string
	switch {
	case userName == "" && password == "":
		msg = "Missing username and password"
	case userName == "":
		msg = "Missing username"
	case password == "":
		msg = "Missing password"
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
		err := app.tmpls.ExecuteTemplate(w, "login.html", err)
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
			log.Printf("Username %q does not exist", userName)
			return sessionToken, sessionExpires, errors.New("login failed")
		}
		log.Println("Login failed", err)
		return sessionToken, sessionExpires, errors.New("login failed")
	}

	// compared hashed password with given password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password for %q", userName)
		return sessionToken, sessionExpires, errors.New("login failed")
	}

	// create a new random sessions token
	sessionToken, err = GenerateRandomString(32)
	if err != nil {
		log.Print("Could not generate sessionToken")
		log.Print(err)
		return sessionToken, sessionExpires, errors.New("login failed")
	}

	// store the sessionToken
	sessionExpires = time.Now().Add(time.Duration(app.config.SessionExpiresHours) * time.Hour)

	_, err = app.db.Query("UPDATE users SET sessionToken = ?, sessionExpires = ? WHERE username = ?", sessionToken, sessionExpires, userName)
	if err != nil {
		log.Printf("Unable to store sessionToken for %q", userName)
		log.Print(err)
		return sessionToken, sessionExpires, errors.New("login failed")
	}

	return sessionToken, sessionExpires, nil
}
