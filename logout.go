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
	"log"
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
		log.Println("invalid method", r.Method)
		return
	}

	sessionTokenValue, err := GetCookieValue(r, "sessionToken")
	if err != nil {
		log.Println("error getting session token cookie", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// create an empty sessionToken cookie with negative MaxAge to delete
	http.SetCookie(w, &http.Cookie{Name: "sessionToken", Value: "", MaxAge: -1})

	// remove session from database
	// TODO: consider removing all sessions for user
	if sessionTokenValue != "" {
		err := RemoveToken(app.DB, "session", sessionTokenValue)
		if err != nil {
			log.Printf("remove token failed for %q: %v", sessionTokenValue, err)
			// TODO: display error or just continue?
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	// display page
	err = RenderTemplate(app.Tmpls, w, "logout.html",
		LogoutPageData{Title: app.Config.Title})
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}
