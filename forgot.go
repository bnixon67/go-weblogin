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

// ForgotPageData contains data passed to the HTML template.
type ForgotPageData struct {
	Title     string
	Message   string
	EmailFrom string
}

// ForgotHandler handles /forgot requests.
func (app *App) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	// only allow valid methods
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		log.Println("invalid method", r.Method)
		return
	}

	// dispatch based on method
	switch r.Method {

	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "forgot.html",
			ForgotPageData{Title: app.Config.Title})
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}

	case http.MethodPost:
		app.forgotPost(w, r)

	}
}

const (
	MsgMissingEmail  = "Please provide an Email."
	MsgNoSuchUser    = "There is no user for the Email provided."
	MsgMissingAction = "Action is missing."
	MsgInvalidAction = "Action is invalid."
)

// forgotPost is called for the POST method of the ForgotHandler.
func (app *App) forgotPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	email := strings.TrimSpace(r.PostFormValue("email"))
	action := strings.TrimSpace(r.PostFormValue("action"))

	// check for missing values
	var msg string
	switch {
	case action == "":
		msg = MsgMissingAction
	case email == "":
		msg = MsgMissingEmail
	}

	// check for invalid action
	if action != "" {
		allowed := []string{"user", "password"}
		if !StringContains(allowed, action) {
			msg = MsgInvalidAction
		}
	}

	// if error msg, display and return
	if msg != "" {
		log.Println(msg)
		pageData := ForgotPageData{
			Title: app.Config.Title, Message: msg,
		}
		err := RenderTemplate(app.Tmpls, w, "forgot.html", pageData)
		if err != nil {
			log.Printf("error executing template: %v", err)
			return
		}
		return
	}

	// get userName for email if provided on the form
	var userName string
	if email != "" {
		var err error
		userName, err = GetUserNameForEmail(app.DB, email)
		if err != nil || userName == "" {
			log.Printf("failed to GetUserNameForEmail %q: %v",
				email, err)
			msg = MsgNoSuchUser
		}
	}

	var emailText string
	switch {
	case userName == "":
		emailText = fmt.Sprintf("This email address is not registered for %s.", app.Config.Title)

	case action == "password":
		// create and save a new session token
		// TODO: use config value for ResetExpiresHours
		resetToken, err := SaveNewToken(app.DB, "reset", userName, 12, 1)
		if err != nil {
			log.Printf("unable to save reset token: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		emailText = fmt.Sprintf("Please vist https://%s/reset?rtoken=%s to reset your password for %s", app.Config.BaseURL, resetToken.Value, app.Config.Title)

	case action == "user":
		emailText = fmt.Sprintf("Your User Name is %s for %s", userName, app.Config.Title)
	}

	err := SendEmail(app.Config.SMTPUser, app.Config.SMTPPassword,
		app.Config.SMTPHost, app.Config.SMTPPort, email,
		app.Config.Title, emailText)
	if err != nil {
		log.Printf("unable to SendEmail: %v", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	err = RenderTemplate(app.Tmpls, w, "forgot_sent.html",
		ForgotPageData{
			Title:     app.Config.Title,
			EmailFrom: app.Config.SMTPUser,
		})
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}
