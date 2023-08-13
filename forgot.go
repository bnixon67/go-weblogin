/*
Copyright 2023 Bill Nixon

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
	"log/slog"
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
	logger := slog.With(slog.Group("request",
		slog.String("id", GetReqID(r.Context())),
		slog.String("remoteAddr", GetRealRemoteAddr(r)),
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
	))

	// only allow valid methods
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		logger.Error("invalid HTTP method")
		return
	}

	// dispatch based on method
	switch r.Method {

	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "forgot.html",
			ForgotPageData{Title: app.Config.Title})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		logger.Info("ForgotHandler")

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
	logger := slog.With(slog.Group("request",
		slog.String("id", GetReqID(r.Context())),
		slog.String("remoteAddr", GetRealRemoteAddr(r)),
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
	))

	// get form values
	email := strings.TrimSpace(r.PostFormValue("email"))
	action := strings.TrimSpace(r.PostFormValue("action"))

	logger.Info("forgotPost", "email", email, "action", action)

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
		logger.Warn("error", "display", msg)
		pageData := ForgotPageData{
			Title: app.Config.Title, Message: msg,
		}
		err := RenderTemplate(app.Tmpls, w, "forgot.html", pageData)
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
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
			logger.Error("failed to GetUserNameForEmail",
				"email", email,
				"err", err)
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
			logger.Error("unable to save reset token", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		emailText = fmt.Sprintf("Please vist %s/reset?rtoken=%s to reset your password for %s", app.Config.BaseURL, resetToken.Value, app.Config.Title)

	case action == "user":
		emailText = fmt.Sprintf("Your User Name is %s for %s", userName, app.Config.Title)
	}

	subj := app.Config.Title + " " + action
	err := SendEmail(app.Config.SMTP.User, app.Config.SMTP.Password, app.Config.SMTP.Host, app.Config.SMTP.Port, email, subj, emailText)
	if err != nil {
		logger.Error("unable to SendEmail", "err", err)
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	logger.Info("sent email",
		slog.Group("email",
			slog.String("to", email),
			slog.String("subject", subj),
			slog.String("emailText", emailText),
		),
	)

	err = RenderTemplate(app.Tmpls, w, "forgot_sent.html",
		ForgotPageData{
			Title:     app.Config.Title,
			EmailFrom: app.Config.SMTP.User,
		})
	if err != nil {
		logger.Error("unable to RenderTemplate", "err", err)
		return
	}
}
