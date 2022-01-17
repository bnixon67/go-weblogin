package main

import (
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

// MustOrHTTPError is a helper that wraps a call to a function returning err and uses w to call http.Error to return an InternalServerError.
func MustOrHTTPError(w http.ResponseWriter, err error) error {
	code := http.StatusInternalServerError

	if err != nil {
		http.Error(w, http.StatusText(code), code)
	}
	return err
}
