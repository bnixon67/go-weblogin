package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles /register requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RegisterHandler", r.Method)

	switch r.Method {

	case "GET":
		err := tmpls.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "POST":
		registerPut(w, r)

	default:
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// registerPut is called for the PUT method of the RegisterHandler
func registerPut(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := r.PostFormValue("userName")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	email := r.PostFormValue("email")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	// check for missing values
	// redundant given client side required fields, but good practice
	if userName == "" || password1 == "" || password2 == "" || firstName == "" || lastName == "" || email == "" {
		msg := "Missing values"
		log.Println(msg, "for", userName)
		err := tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check that userName doesn't already exist
	exists, _ := CheckForUserName(userName)
	if exists {
		log.Printf("UserName exists for %q", userName)
		err := tmpls.ExecuteTemplate(w, "register.html", "Sorry, your desired User Name already exists. Please try a different User Name")
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := "Passwords do not match"
		log.Println(msg, "for", userName)
		err := tmpls.ExecuteTemplate(w, "register.html", msg)
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
		err := tmpls.ExecuteTemplate(w, "register.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// store the user and hashed password
	_, err = db.Query("INSERT INTO users(username, hashedPassword, firstName, lastName, email) VALUES (?, ?, ?, ?, ?)",
		userName, string(hashedPassword), firstName, lastName, email)
	if err != nil {
		msg := "Unable to register user"
		log.Println(msg, "for", userName, err)
		err := tmpls.ExecuteTemplate(w, "register.html", msg)
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
