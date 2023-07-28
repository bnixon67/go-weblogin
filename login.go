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
	"net/http"
	"strings"

	"golang.org/x/exp/slog"
)

// LoginPageData contains data passed to the HTML template.
type LoginPageData struct {
	Title   string
	Message string
}

// LoginHandler handles /login requests.
func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet, http.MethodPost}) {
		slog.Error("invalid HTTP method", "method", r.Method)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{Title: app.Config.Title})
		if err != nil {
			slog.Error("unabel to RenderTemplate", "err", err)
			return
		}
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
		slog.Info("error", "display", msg)
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{Title: app.Config.Title, Message: msg})
		if err != nil {
			slog.Error("uanble to RenderTemplate", "err", err)
			return
		}
		return
	}

	// attempt to login the given userName with the given password
	token, err := app.LoginUser(userName, password)
	if err != nil {
		slog.Error("failed to LoginUser", "userName", userName, "err", err)
		err := RenderTemplate(app.Tmpls, w, "login.html",
			LoginPageData{
				Title:   app.Config.Title,
				Message: MsgLoginFailed,
			})
		if err != nil {
			slog.Error("unable to RenderTemplate", "err", err)
			return
		}
		return
	}

	// login successful, so create a cookie for the session Token
	slog.Info("successful login", "userName", userName)
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
}

// LoginUser returns a session Token if userName and password is correct.
func (app *App) LoginUser(userName, password string) (Token, error) {
	err := CompareUserPassword(app.DB, userName, password)
	if err != nil {
		WriteEvent(app.DB, Event{UserName: userName, Action: ActionLogin, Result: false, Message: err.Error()})
		return Token{}, err
	}

	// create and save a new session token
	token, err := SaveNewToken(app.DB, "session", userName, 32, app.Config.SessionExpiresHours)
	if err != nil {
		WriteEvent(app.DB, Event{UserName: userName, Action: ActionSaveToken, Result: false})
		slog.Error("unable to SaveNewToken", "err", err)
		return Token{}, fmt.Errorf("unable to save token: %w", err)
	}

	WriteEvent(app.DB, Event{UserName: userName, Action: ActionLogin, Result: true})
	return token, nil
}
