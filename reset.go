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

	"golang.org/x/crypto/bcrypt"
)

// ResetPageData contains data passed to the HTML template.
type ResetPageData struct {
	Title      string
	Message    string
	ResetToken string
}

// ResetHandler handles /reset requests.
func (app *App) ResetHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.With(slog.Group("request",
		slog.String("id", GetReqID(r.Context())),
		slog.String("remoteAddr", GetRealRemoteAddr(r)),
		slog.String("method", r.Method),
		slog.String("url", r.RequestURI),
	))

	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		logger.Error("invalid HTTP method", "method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "reset.html",
			ResetPageData{
				Title:      app.Cfg.Title,
				ResetToken: r.URL.Query().Get("rtoken"),
			})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		logger.Info("ResetHandler")

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

	logger := slog.With(
		slog.Group("request",
			slog.String("id", GetReqID(r.Context())),
			slog.String("remoteAddr", GetRealRemoteAddr(r)),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
		),
	)

	// check for missing values
	// redundant given client side required fields, but good practice
	if resetToken == "" || password1 == "" || password2 == "" {
		msg := MsgMissingRequired
		logger.Warn("missing field(s)",
			slog.Group("form",
				"rtoken empty", resetToken == "",
				"password1 empty", password1 == "",
				"password2 empty", password2 == "",
			),
		)
		err := RenderTemplate(app.Tmpls, w, tmplFileName,
			ResetPageData{
				Title:      app.Cfg.Title,
				Message:    msg,
				ResetToken: r.URL.Query().Get("rtoken"),
			})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	// check that password fields match
	// may be redundant if done client side, but good practice
	if password1 != password2 {
		msg := MsgPasswordsDifferent
		logger.Warn("passwords don't match")
		err := RenderTemplate(app.Tmpls, w, tmplFileName,
			ResetPageData{
				Title:      app.Cfg.Title,
				Message:    msg,
				ResetToken: r.URL.Query().Get("rtoken"),
			})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	userName, err := GetUserNameForResetToken(app.DB, resetToken)
	if err != nil {
		logger.Error("failed GetUserNameForResetToken",
			"resetToken", resetToken,
			"err", err)
		msg := "Please provide a valid Reset Token"
		err := RenderTemplate(app.Tmpls, w, tmplFileName,
			ResetPageData{
				Title:      app.Cfg.Title,
				Message:    msg,
				ResetToken: r.URL.Query().Get("rtoken"),
			})
		if err != nil {
			logger.Error("failed to RenderTemplate", "err", err)
			return
		}
		return
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		msg := "Cannot hash password"
		logger.Error("failed bcrypt.GenerateFromPassword",
			"userName", userName, "err", err)
		err := RenderTemplate(app.Tmpls, w, tmplFileName,
			ResetPageData{Title: app.Cfg.Title, Message: msg})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	// store the user and hashed password
	_, err = app.DB.Exec("UPDATE users SET hashedPassword = ? WHERE username = ?", string(hashedPassword), userName)
	if err != nil {
		logger.Error("update password failed",
			"userName", userName, "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// TODO: don't allow reuse of the reset token if successful

	// register successful
	logger.Info("successful password reset", "userName", userName)
	WriteEvent(app.DB, EventReset, true, userName, "success")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
