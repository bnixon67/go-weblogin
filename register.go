package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// TODO: handle duplicate emails

// RegisterHandler handles /register requests
func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := app.tmpls.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodPost:
		app.registerPost(w, r)
	}
}

const (
	MSG_REGISTER_MISSING_VALUES           = "Please provide all the required values"
	MSG_REGISTER_ERR_EXISTING_USER        = "Sorry, your desired User Name already exists. Please try a different User Name"
	MSG_REGISTER_ERR_MISMATCHED_PASSWORDS = "Password do not match."
)

// registerPost is called for the POST method of the RegisterHandler
func (app *App) registerPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := r.PostFormValue("userName")
	fullName := r.PostFormValue("fullName")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	// check for missing values
	// redundant given client side required fields, but good practice
	if userName == "" || password1 == "" || password2 == "" || fullName == "" || email == "" {
		msg := MSG_REGISTER_MISSING_VALUES
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check that userName doesn't already exist
	exists, _ := app.CheckForUserName(userName)
	if exists {
		log.Printf("UserName exists for %q", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", MSG_REGISTER_ERR_EXISTING_USER)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MSG_REGISTER_ERR_MISMATCHED_PASSWORDS
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		msg := "Cannot hash password"
		log.Println(msg, "for", userName)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// store the user and hashed password
	_, err = app.db.Query("INSERT INTO users(username, hashedPassword, fullName, email) VALUES (?, ?, ?, ?)",
		userName, string(hashedPassword), fullName, email)
	if err != nil {
		msg := "Unable to register user"
		log.Println(msg, "for", userName, err)
		err := app.tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// register successful
	log.Printf("Username %q registered", userName)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
