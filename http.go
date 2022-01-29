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
	"html/template"
	"log"
	"net/http"
	"strings"
)

// StringContains reports whether val is within arr.
func StringContains(arr []string, val string) bool {
	for _, e := range arr {
		if e == val {
			return true
		}
	}
	return false
}

// ValidMethod reports if r.Method is within allowed. If r.Method is not allowed or is OPTIONS, then w is updated with an appropriate response, false is returned, and any Handler using this function should return.
func ValidMethod(w http.ResponseWriter, r *http.Request, allowed []string) bool {
	if StringContains(allowed, r.Method) {
		return true
	}

	allowed = append(allowed, http.MethodOptions)
	w.Header().Set("Allow", strings.Join(allowed, ", "))

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	return false
}

const MsgTemplateError = "Sorry, the server was unable to display this page. Please contact the administrator."

// RenderTemplate is a helper to call template.ExecuteTemplate and returns a http.Error unpon failure. Like http.Error, it does not otherwise end the request, so the caller must ensure no further writes are done to w if non-nil is returned.
func RenderTemplate(t *template.Template, w http.ResponseWriter, name string, data interface{}) error {
	err := t.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, MsgTemplateError, http.StatusInternalServerError)
	}

	return err
}

// ServeFileHandler is a simple http.ServeFile wrapper.
func ServeFileHandler(file string) http.HandlerFunc {
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

	log.Println(ip, r.Method, r.RequestURI)

	l.Next.ServeHTTP(w, r)
}
