package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginPageData record
type LoginPageData struct {
	Message     string
	CurrentUser User
}

// LoginHandler handles /login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("LoginHandler", r.Method)

	switch r.Method {

	case "GET":
		// get currentUser if sessionToken exists
		var sessionToken string
		var currentUser User
		var err error
		c, err := r.Cookie("sessionToken")
		if err == nil {
			sessionToken = c.Value
			currentUser, _ = GetUserForSessionToken(sessionToken)
			log.Printf("%+v", currentUser)
		}

		err = tmpls.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Println("Error in executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "POST":
		loginPut(w, r)

	default:
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// loginPut is called for the PUT method of the LoginHandler
func loginPut(w http.ResponseWriter, r *http.Request) {
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
		err := tmpls.ExecuteTemplate(w, "login.html", msg)
		if err != nil {
			log.Println("Error in executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// attempt to login the given userName with the given password
	sessionToken, sessionExpires, err := loginUser(userName, password)
	if err != nil {
		err := tmpls.ExecuteTemplate(w, "login.html", err)
		if err != nil {
			log.Println("Error in executing template", err)
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
func loginUser(userName, password string) (string, time.Time, error) {
	var sessionToken string
	var sessionExpires time.Time

	// get hashed password for the given user
	result := db.QueryRow("SELECT hashedPassword FROM users WHERE username=?", userName)
	var hashedPassword string
	err := result.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Username %q does not exist", userName)
			return sessionToken, sessionExpires, errors.New("Login failed")
		}
		log.Println("Login failed", err)
		return sessionToken, sessionExpires, errors.New("Login failed")
	}

	// compared hashed password with given password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password for %q", userName)
		return sessionToken, sessionExpires, errors.New("Login failed")
	}

	// create a new random sessions token
	sessionToken = uuid.NewString()

	// store the sessionToken
	sessionExpires = time.Now().Add(120 * time.Second)

	_, err = db.Query("UPDATE users SET sessionToken = ?, sessionExpires = ? WHERE username = ?", sessionToken, sessionExpires, userName)
	if err != nil {
		log.Printf("Unable to store sessionToken for %q", userName)
		log.Print(err)
		return sessionToken, sessionExpires, errors.New("Login failed")
	}

	return sessionToken, sessionExpires, nil
}
