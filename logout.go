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
)

// LogoutPageData contains data passed to the HTML template.
type LogoutPageData struct {
	Title   string
	Message string
}

// LogoutHandler handles /logout requests.
func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet}) {
		slog.Error("invalid HTTP method", "method", r.Method)
		return
	}

	user, err := GetUser(w, r, app.DB)
	if err != nil {
		slog.Error("failed to GetUser", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("logout", "user", user)
	if user.UserName != "" {
		WriteEvent(app.DB, EventLogout, true, user.UserName, "")
	}

	sessionTokenValue, err := GetCookieValue(r, SessionTokenCookieName)
	if err != nil {
		slog.Error("failed to GetCookieValue", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create an empty sessionToken cookie with negative MaxAge to delete
	http.SetCookie(w,
		&http.Cookie{
			Name: SessionTokenCookieName, Value: "", MaxAge: -1,
		})

	// remove session from database
	// TODO: consider removing all sessions for user
	if sessionTokenValue != "" {
		err := RemoveToken(app.DB, "session", sessionTokenValue)
		if err != nil {
			slog.Error("filed to RemoveToken",
				"sessionTokenValue", sessionTokenValue,
				"err", err)
			// TODO: display error or just continue?
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// display page
	err = RenderTemplate(app.Tmpls, w, "logout.html",
		LogoutPageData{Title: app.Config.Title})
	if err != nil {
		slog.Error("failed to RenderTemplate", "err", err)
		return
	}
}
