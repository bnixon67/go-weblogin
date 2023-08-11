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
	"log/slog"
	"net/http"
	"strings"
)

const (
	MsgMissingRequired    = "Please provide all the required values"
	MsgUserNameExists     = "User Name already exists."
	MsgEmailExists        = "Email Address already registered."
	MsgPasswordsDifferent = "Password values do not match."
)

// RegisterPageData contains data passed to the HTML template.
type RegisterPageData struct {
	Title   string
	Message string
}

// RegisterHandler handles /register requests.
func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		slog.Error("invalid HTTP method", "method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{Title: app.Config.Title})
		if err != nil {
			slog.Error("unable to parse template", "err", err)
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
	if IsEmpty(userName, fullName, email, password1, password2) {
		msg := MsgMissingRequired
		slog.Warn("register missing values",
			"userName", userName,
			"fullName", fullName,
			"email", email,
			"password1 empty", password1 == "",
			"password2 empty", password2 == "")
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title: app.Config.Title, Message: msg,
			})
		if err != nil {
			slog.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that password fields match
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		slog.Warn("passwords do not match", "userName", userName)
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title: app.Config.Title, Message: msg,
			})
		if err != nil {
			slog.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that userName doesn't already exist
	userExists, err := UserExists(app.DB, userName)
	if err != nil {
		slog.Error("UserExists failed",
			"userName", userName, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		slog.Warn("user exists", "userName", userName)
		WriteEvent(app.DB, Event{UserName: userName, Action: ActionRegister, Result: false, Message: "user exists"})
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Config.Title,
				Message: MsgUserNameExists,
			})
		if err != nil {
			slog.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that email doesn't already exist
	emailExists, err := EmailExists(app.DB, email)
	if err != nil {
		slog.Error("EmailExists failed", "email", email, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if emailExists {
		slog.Warn("email exists", "email", email)
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Config.Title,
				Message: MsgEmailExists,
			})
		if err != nil {
			slog.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// Register User
	err = RegisterUser(app.DB, userName, fullName, email, password1)
	if err != nil {
		slog.Error("RegisterUser failed",
			"userName", userName, "err", err)
		WriteEvent(app.DB, Event{UserName: userName, Action: ActionRegister, Result: false, Message: err.Error()})
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Config.Title,
				Message: "Unable to Register User",
			})
		if err != nil {
			slog.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// registration successful
	slog.Info("registered user", "userName", userName)
	WriteEvent(app.DB, Event{UserName: userName, Action: ActionRegister, Result: true})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
