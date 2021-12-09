package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogoutHandlerInvalidMethod(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)

	app.LogoutHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func getCookie(name string, cookies []*http.Cookie) (*http.Cookie, error) {
	for _, c := range cookies {
		if name == c.Name {
			return c, nil
		}
	}
	return nil, http.ErrNoCookie
}

func TestLogoutHandlerGet(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/logout", nil)

	app.LogoutHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "You have been logged out."
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}

	c, err := getCookie("sessionToken", w.Result().Cookies())
	if err != nil {
		t.Fatalf("sessionToken cookie missing")
	}
	if c.Value != "" {
		t.Fatalf("sessionToken not empty")
	}
	if c.MaxAge != -1 {
		t.Fatalf("sessionToken.MaxAge not -1")
	}
}
