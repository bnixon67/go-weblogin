/*
Copyright 2022 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// ForgotHandler handles /forgot requests.
func (app *App) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "forgot.html", nil)
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
		err := RenderTemplate(app.Tmpls, w, "forgot.html", MsgMissingEmail)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// get userName for email provided on the form
	// TODO: send email rather than exposing no user for email
	userName, err := GetUserNameForEmail(app.DB, email)
	if err != nil || userName == "" {
		log.Printf("failed to GetUserNameForEmail %q: %v", email, err)
		err := RenderTemplate(app.Tmpls, w, "forgot.html", MsgNoSuchUser)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// create and save a new session token
	// TODO: use config value for ResetExpiresHours
	resetToken, err := SaveNewToken(app.DB, "reset", userName, 12, 1)
	if err != nil {
		log.Printf("unable to save reset token: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resetURL := fmt.Sprintf("https://%s:%s/reset?rtoken=%s", app.Config.ServerHost, app.Config.ServerPort, resetToken.Value)
	err = SendEmail(app.Config.SMTPUser, app.Config.SMTPPassword, app.Config.SMTPHost, app.Config.SMTPPort, email, "Reset Pasword", resetURL)
	if err != nil {
		log.Printf("unable to SendEmail: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = RenderTemplate(app.Tmpls, w, "forgot_sent.html", nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}
