package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ForgotHandler handles /forgot requests
func (app *App) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case http.MethodPost:
		app.forgotPost(w, r)
	}
}

const (
	MSG_MISSING_EMAIL = "Please provide an Email"
	MSG_NO_SUCH_USER  = "There is no registered User Name for the Email provided."
)

// forgotPost is called for the POST method of the LoginHandler
func (app *App) forgotPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	email := strings.TrimSpace(r.PostFormValue("email"))

	// check for missing values
	if email == "" {
		log.Print("email is empty")
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", MSG_MISSING_EMAIL)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// get userName for email provided on the form
	// TODO: send email rather than exposing no user for email
	userName, err := GetUserNameForEmail(app.db, email)
	if err != nil || userName == "" {
		log.Printf("failed to GetUserNameForEmail %q: %v", email, err)
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", MSG_NO_SUCH_USER)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// create and save a new session token
	// TODO: use config value for ResetExpiresHours
	resetToken, err := SaveNewToken(app.db, "reset", userName, 12, 1)
	if err != nil {
		log.Printf("unable to save reset token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resetURL := fmt.Sprintf("https://%s:%s/reset?rtoken=%s", app.config.ServerHost, app.config.ServerPort, resetToken.Value)
	SendEmail(app.config.SmtpUser, app.config.SmtpPassword, app.config.SmtpHost, app.config.SmtpPort, email, "Reset Pasword", resetURL)

	err = app.tmpls.ExecuteTemplate(w, "forgot_sent.html", err)
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
