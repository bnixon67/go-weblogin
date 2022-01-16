package main

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type ResetData struct {
	Msg        string
	ResetToken string
}

// ResetHandler handles /rest requests.
func (app *App) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := app.tmpls.ExecuteTemplate(w, "reset.html", ResetData{ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		app.resetPost(w, r, "reset.html")
	}
}

// resetPost is called for the POST method of the RegisterHandler.
func (app *App) resetPost(w http.ResponseWriter, r *http.Request, tmplFileName string) {
	// get form values
	resetToken := strings.TrimSpace(r.PostFormValue("rtoken"))
	password1 := strings.TrimSpace(r.PostFormValue("password1"))
	password2 := strings.TrimSpace(r.PostFormValue("password2"))

	// check for missing values
	// redundant given client side required fields, but good practice
	if resetToken == "" || password1 == "" || password2 == "" {
		msg := MsgMissingRequired
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	userName, err := GetUserNameForResetToken(app.db, resetToken)
	if err != nil {
		log.Printf("invalid reset token: %v", err)
		msg := "Please provide a valid Reset Token"
		err := app.tmpls.ExecuteTemplate(w, tmplFileName, ResetData{Msg: msg, ResetToken: r.URL.Query().Get("rtoken")})
		if err != nil {
			log.Println("error executing template", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	// store the user and hashed password
	_, err = app.db.Exec("UPDATE users SET hashedPassword = ? WHERE username = ?", string(hashedPassword), userName)
	if err != nil {
		log.Printf("update password failed for %q: %v", userName, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// register successful
	log.Printf("Password reset for %q", userName)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
