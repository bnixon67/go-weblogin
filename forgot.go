package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// ForgotHandler handles /forgot requests.
func (app *App) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.tmpls, w, "forgot.html", nil)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}

	case http.MethodPost:
		app.forgotPost(w, r)
	}
}

const (
	MsgMissingEmail = "Please provide an Email"
	MsgNoSuchUser   = "There is no registered User Name for the Email provided."
)

// forgotPost is called for the POST method of the LoginHandler.
func (app *App) forgotPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	email := strings.TrimSpace(r.PostFormValue("email"))

	// check for missing values
	if email == "" {
		log.Print("email is empty")
		err := RenderTemplate(app.tmpls, w, "forgot.html", MsgMissingEmail)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// get userName for email provided on the form
	// TODO: send email rather than exposing no user for email
	userName, err := GetUserNameForEmail(app.db, email)
	if err != nil || userName == "" {
		log.Printf("failed to GetUserNameForEmail %q: %v", email, err)
		err := RenderTemplate(app.tmpls, w, "forgot.html", MsgNoSuchUser)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// create and save a new session token
	// TODO: use config value for ResetExpiresHours
	resetToken, err := SaveNewToken(app.db, "reset", userName, 12, 1)
	if err != nil {
		log.Printf("unable to save reset token: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resetURL := fmt.Sprintf("https://%s:%s/reset?rtoken=%s", app.config.ServerHost, app.config.ServerPort, resetToken.Value)
	err = SendEmail(app.config.SMTPUser, app.config.SMTPPassword, app.config.SMTPHost, app.config.SMTPPort, email, "Reset Pasword", resetURL)
	if err != nil {
		log.Printf("unable to SendEmail: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = RenderTemplate(app.tmpls, w, "forgot_sent.html", nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}
