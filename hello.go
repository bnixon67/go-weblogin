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
)

// HelloPageData contains data passed to the HTML template.
type HelloPageData struct {
	Title   string
	Message string
	User    User
}

// HelloHandler prints a simple hello message.
func (app *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet}) {
		slog.Error("invalid HTTP method", "method", r.Method)
		return
	}

	currentUser, err := GetUser(w, r, app.DB)
	if err != nil {
		slog.Error("failed to GetUser", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// display page
	err = RenderTemplate(app.Tmpls, w, "hello.html", HelloPageData{Message: "", User: currentUser})
	if err != nil {
		slog.Error("unable to RenderTemplate", "err", err)
		return
	}
}
