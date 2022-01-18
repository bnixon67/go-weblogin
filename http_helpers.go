package main

import (
	"html/template"
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
