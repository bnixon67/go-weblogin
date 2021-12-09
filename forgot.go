package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// ForgotHandler handles /forgot requests
func (app *App) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)

	switch r.Method {

	case "GET":
		err := app.tmpls.ExecuteTemplate(w, "forgot.html", nil)
		if err != nil {
			log.Println("error executing template", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	case "POST":
		app.forgotPut(w, r)

	default:
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

const (
	MSG_MISSING_EMAIL = "Missing email"
	MSG_NO_SUCH_USER  = "There is no registered User Name for the Email provided."
)

// forgotPut is called for the PUT method of the LoginHandler
func (app *App) forgotPut(w http.ResponseWriter, r *http.Request) {
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
