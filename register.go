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
	logger := slog.With(slog.Group("request",
		slog.String("id", GetReqID(r.Context())),
		slog.String("remoteAddr", GetRealRemoteAddr(r)),
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
	))

	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		logger.Error("invalid HTTP method")
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{Title: app.Cfg.Title})
		if err != nil {
			logger.Error("unable to parse template", "err", err)
			return
		}
		logger.Info("RegisterHandler")

	case http.MethodPost:
		app.registerPost(w, r)
	}
}

// registerPost is called for the POST method of the RegisterHandler.
func (app *App) registerPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("userName"))
	fullName := strings.TrimSpace(r.PostFormValue("fullName"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	password1 := strings.TrimSpace(r.PostFormValue("password1"))
	password2 := strings.TrimSpace(r.PostFormValue("password2"))

	logger := slog.With(
		slog.Group("request",
			slog.String("id", GetReqID(r.Context())),
			slog.String("remoteAddr", GetRealRemoteAddr(r)),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
		),
		slog.Group("form",
			"userName", userName,
			"fullName", fullName,
			"email", email,
			"password1 empty", password1 == "",
			"password2 empty", password2 == "",
		),
	)

	// check for missing values
	if IsEmpty(userName, fullName, email, password1, password2) {
		msg := MsgMissingRequired
		logger.Warn("missing values")
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title: app.Cfg.Title, Message: msg,
			})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that password fields match
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		logger.Warn("passwords do not match")
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title: app.Cfg.Title, Message: msg,
			})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that userName doesn't already exist
	userExists, err := UserExists(app.DB, userName)
	if err != nil {
		logger.Error("UserExists failed", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if userExists {
		logger.Warn("user already exists")
		WriteEvent(app.DB, EventRegister, false, userName, "user already exists")
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Cfg.Title,
				Message: MsgUserNameExists,
			})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// check that email doesn't already exist
	emailExists, err := EmailExists(app.DB, email)
	if err != nil {
		logger.Error("EmailExists failed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if emailExists {
		logger.Warn("email already exists")
		WriteEvent(app.DB, EventRegister, false, userName, "email already exists")
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Cfg.Title,
				Message: MsgEmailExists,
			})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// Register User
	err = RegisterUser(app.DB, userName, fullName, email, password1)
	if err != nil {
		logger.Error("RegisterUser failed", "err", err)
		WriteEvent(app.DB, EventRegister, false, userName, err.Error())
		err := RenderTemplate(app.Tmpls, w, "register.html",
			RegisterPageData{
				Title:   app.Cfg.Title,
				Message: "Unable to Register User",
			})
		if err != nil {
			logger.Error("unable to execute template", "err", err)
			return
		}
		return
	}

	// registration successful
	logger.Info("registered user")
	WriteEvent(app.DB, EventRegister, true, userName, "success")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
