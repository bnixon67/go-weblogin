package main

import (
	"log"
	"net/http"
	"strings"
)

const (
	MsgMissingRequired    = "Please provide all the required values"
	MsgUserNameExists     = "Your desired User Name already exists."
	MsgEmailExists        = "A User Name already exists for this Email Address."
	MsgPasswordsDifferent = "Password do not match."
)

// RegisterHandler handles /register requests.
func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", nil))
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}

	case http.MethodPost:
		app.registerPost(w, r)
	}
}

// IsEmpty returns true if any of the strings are empty, otherwise false.
func IsEmpty(strs ...string) bool {
	for _, s := range strs {
		if s == "" {
			return true
		}
	}

	return false
}

// registerPost is called for the POST method of the RegisterHandler.
func (app *App) registerPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("userName"))
	fullName := strings.TrimSpace(r.PostFormValue("fullName"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password1 := strings.TrimSpace(r.PostFormValue("password1"))
	password2 := strings.TrimSpace(r.PostFormValue("password2"))

	// check for missing values
	// redundant given client side required fields, but good practice
	if IsEmpty(userName, fullName, email, password1, password2) {
		msg := MsgMissingRequired
		log.Println(msg, "for", userName)
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", msg))
		if err != nil {
			log.Println("error executing template", err)
			return
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		log.Println(msg, "for", userName)
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", msg))
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// check that userName doesn't already exist
	userExists, err := UserExists(app.db, userName)
	if err != nil {
		log.Printf("error in UserExists for %q: %v", userName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		log.Printf("userName %q already exists", userName)
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", MsgUserNameExists))
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// check that email doesn't already exist
	emailExists, err := EmailExists(app.db, email)
	if err != nil {
		log.Printf("error in EmailExists for %q: %v", email, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if emailExists {
		log.Printf("email %q already exists", email)
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", MsgEmailExists))
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// Register User
	err = RegisterUser(app.db, userName, fullName, email, password1)
	if err != nil {
		log.Printf("unable to RegisterUser %q: %v", userName, err)
		err := MustOrHTTPError(w, app.tmpls.ExecuteTemplate(w, "register.html", "Unable to Register User"))
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// registration successful
	log.Printf("Username %q registered", userName)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
