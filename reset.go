package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type ResetData struct {
	Msg        string
	ResetToken string
}

// ResetHandler handles /rest requests
func (app *App) ResetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "from", r.RemoteAddr)

	switch r.Method {

	case "GET":
		err := app.tmpls.ExecuteTemplate(w, "reset.html", ResetData{ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "POST":
		app.resetPut(w, r, "reset.html")

	default:
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

const (
	MSG_MISSING_RESET_VALUES = "Missing values"
	MSG_MISMATCHED_PASSWORDS = "Passwords do not match"
)

// resetPut is called for the PUT method of the RegisterHandler
func (app *App) resetPut(w http.ResponseWriter, r *http.Request, tmplFileName string) {
	// get form values
	resetToken := r.PostFormValue("rtoken")
	password1 := r.PostFormValue("password1")
	password2 := r.PostFormValue("password2")

	// check for missing values
	// redundant given client side required fields, but good practice
	if resetToken == "" || password1 == "" || password2 == "" {
		msg := MSG_MISSING_RESET_VALUES
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MSG_MISMATCHED_PASSWORDS
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	userName, err := app.GetUserNameForResetToken(resetToken)
	if err != nil {
		log.Println("Invalid Reset Token")
		msg := "Please provide a valid Reset Token"
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
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
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// store the user and hashed password
	_, err = app.db.Query("UPDATE users SET hashedPassword = ? WHERE username = ?", string(hashedPassword), userName)
	if err != nil {
		log.Printf("Unable to update password for %q", userName)
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// register successful
	log.Printf("Password reset for %q", userName)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
