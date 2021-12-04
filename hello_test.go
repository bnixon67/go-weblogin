package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelloHandlerInvalidMethod(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/hello", nil)

	app.HelloHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestHelloHandlerWithoutCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)

	app.HelloHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "You must <a href=\"/login\">login</a>"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestHelloHandlerWithBadSessionToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	r.AddCookie(&http.Cookie{Name: "sessionToken", Value: "foo"})

	app.HelloHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "You must <a href=\"/login\">login</a>"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestHelloHandlerWithGoodSessionToken(t *testing.T) {
	// TODO: better way to define a test user
	token, _, err := app.loginUser("test", "password")
	if err != nil {
		t.Fatalf("could not login user to get session token")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	r.AddCookie(&http.Cookie{Name: "sessionToken", Value: token})

	app.HelloHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "You must <a href=\"/login\">login</a>"
	if strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}
