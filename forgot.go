package main

import (
	"fmt"
	"log"
	"net/http"

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
	MSG_MISSING_EMAIL = "Missing email"
	MSG_NO_SUCH_USER  = "There is no registered User Name for the Email provided."
)

// forgotPost is called for the POST method of the LoginHandler
func (app *App) forgotPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	email := r.PostFormValue("email")

	// check for missing values
	var msg string
	if email == "" {
		msg = MSG_MISSING_EMAIL
	}
	if msg != "" {
		log.Println(msg)
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// get userName for email provided on the form
	userName, err := app.GetUserNameForEmail(email)
	if err != nil || userName == "" {
		msg = MSG_NO_SUCH_USER
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", msg)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// get a new random token to reset password
	resetToken, err := GenerateRandomString(32)
	if err != nil {
		log.Print("Could not generate resetToken")
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	app.SaveResetTokenForUser(userName, resetToken)
	if err != nil {
		log.Print("SaveForgotTokenForUser failed")
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resetURL := fmt.Sprintf("https://%s:%s/reset?rtoken=%s", app.config.ServerHost, app.config.ServerPort, resetToken)
	SendEmail(app.config.SmtpUser, app.config.SmtpPassword, app.config.SmtpHost, app.config.SmtpPort, email, "Reset Pasword", resetURL)

	err = app.tmpls.ExecuteTemplate(w, "forgot_sent.html", err)
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
