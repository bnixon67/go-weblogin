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

// LoginPageData contains data passed to the HTML template.
type LoginPageData struct {
	Title   string
	Message string
}

// LoginHandler handles /login requests.
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{Title: app.Config.Title})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		logger.Info("LoginHandler")

	case http.MethodPost:
		app.loginPost(w, r)
	}
}

const (
	MsgMissingUserNameAndPassword = "Missing username and password"
	MsgMissingUserName            = "Missing username"
	MsgMissingPassword            = "Missing password"
	MsgLoginFailed                = "Login Failed"
)

// loginPost is called for the POST method of the LoginHandler.
func (app *App) loginPost(w http.ResponseWriter, r *http.Request) {
	// get form values
	userName := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))

	logger := slog.With(
		slog.Group("request",
			slog.String("id", GetReqID(r.Context())),
			slog.String("remoteAddr", GetRealRemoteAddr(r)),
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
		),
		slog.Group("form",
			"userName", userName,
			"password empty", password == "",
		),
	)

	// check for missing values
	var msg string
	switch {
	case userName == "" && password == "":
		msg = MsgMissingUserNameAndPassword
	case userName == "":
		msg = MsgMissingUserName
	case password == "":
		msg = MsgMissingPassword
	}
	if msg != "" {
		logger.Info("error", "display", msg)
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{Title: app.Config.Title, Message: msg})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	// attempt to login the given userName with the given password
	token, err := app.LoginUser(userName, password)
	if err != nil {
		logger.Error("failed to LoginUser", "err", err)
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{
				Title:   app.Config.Title,
				Message: MsgLoginFailed,
			})
		if err != nil {
			logger.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	// login successful, so create a cookie for the session Token
	http.SetCookie(w, &http.Cookie{
		Name:     SessionTokenCookieName,
		Value:    token.Value,
		Expires:  token.Expires,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	redirect := r.URL.Query().Get("r")
	if redirect == "" {
		redirect = "/"
	}

	// redirect from login page
	http.Redirect(w, r, redirect, http.StatusSeeOther)

	logger.Info("login successful")
}

// LoginUser returns a session Token if userName and password is correct.
func (app *App) LoginUser(userName, password string) (Token, error) {
	err := CompareUserPassword(app.DB, userName, password)
	if err != nil {
		WriteEvent(app.DB, EventLogin, false, userName, err.Error())

		return Token{}, err
	}

	// create and save a new session token
	token, err := SaveNewToken(app.DB, "session", userName, 32, app.Config.SessionExpiresHours)
	if err != nil {
		WriteEvent(app.DB, EventSaveToken, false, userName, err.Error())
		slog.Error("unable to SaveNewToken", "err", err, "userName", userName)
		return Token{}, fmt.Errorf("unable to save token: %w", err)
	}

	WriteEvent(app.DB, EventLogin, true, userName, "success")

	return token, nil
}
