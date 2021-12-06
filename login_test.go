package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestLoginHandlerInvalidMethod(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/login", nil)

	app.LoginHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestLoginHandlerGet(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/login", nil)

	app.LoginHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "<form method=\"post\""
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestLoginHandlerPostMissingUserNameAndPassword(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", nil)

	app.LoginHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_MISSING_USERNAME_PASSWORD
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestLoginHandlerPostMissingPassword(t *testing.T) {
	data := url.Values{
		"username": {"foo"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.LoginHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_MISSING_PASSWORD
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestLoginHandlerPostMissingUserName(t *testing.T) {
	data := url.Values{
		"password": {"foo"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.LoginHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_MISSING_USERNAME
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestLoginHandlerPostInvalidUserNameAndPassword(t *testing.T) {
	data := url.Values{
		"username": {"foo"},
		"password": {"bar"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.LoginHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_LOGIN_FAILED
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestLoginHandlerPostValidUserNameAndPassword(t *testing.T) {
	data := url.Values{
		"username": {"test"},
		"password": {"password"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.LoginHandler(w, r)

	expectedStatus := http.StatusSeeOther
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expected := ""
	if w.Body.String() != expected {
		t.Fatalf("got body %q, expected %q", w.Body, expected)
	}
}
