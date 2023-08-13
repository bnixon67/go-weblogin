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
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

// StringContains reports if slice contain value.
func StringContains(slice []string, value string) bool {
	for _, e := range slice {
		if e == value {
			return true
		}
	}
	return false
}

// ValidMethod checks if the given HTTP request method is allowed based on the provided list of allowed methods. It returns true if the method is allowed, and false otherwise. If the method is not allowed or is OPTIONS, the function updates the response writer appropriately and returns false.  The calling handler should return without further processing.
func ValidMethod(w http.ResponseWriter, r *http.Request, allowed []string) bool {
	// if method is in allowed list, then return
	if StringContains(allowed, r.Method) {
		return true
	}

	// add OPTIONS if method is not allowed
	allowed = append(allowed, http.MethodOptions)

	// set the "Allow" header to allowed methods
	w.Header().Set("Allow", strings.Join(allowed, ", "))

	// if method is OPTIONS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent) // no content returned
		return false
	}

	// method is not allowed and not OPTIONS
	txt := r.Method + " " + http.StatusText(http.StatusMethodNotAllowed)
	http.Error(w, txt, http.StatusMethodNotAllowed)
	return false
}

const MsgTemplateError = "Sorry, the server was unable to display this page. Please contact the administrator."

// RenderTemplate executes the named template with given data to the writer.
// If an error occurs, writer is updated to indicate a Internal Server Error.
// The caller must ensure no further writes are done for a non-nil error.
func RenderTemplate(t *template.Template, w http.ResponseWriter, name string, data interface{}) error {
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, MsgTemplateError, http.StatusInternalServerError)
	}

	return err
}

// ServeFileHandler returns a HandlerFunc to serve the specified file.
func ServeFileHandler(file string) http.HandlerFunc {
	// check if file exists and is accessible
	_, err := os.Stat(file)
	if err != nil {
		slog.Error("does not exist", "file", file)
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}

// LogRequestHandler is middleware that logs all HTTP requests and
// then calls the next HTTP handler specified.
type LogRequestHandler struct {
	Next http.Handler
}

// ServerHTTP for logRequestHandler log the HTTP request and then
// calls the next HTTP handler specified.
func (l *LogRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get real IP address if using Cloudflare or similar service
	var ip string
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}

	slog.Debug("request", "ip", ip, "method", r.Method, "url", r.RequestURI)

	l.Next.ServeHTTP(w, r)
}
