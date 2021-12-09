package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestForgotHandlerInvalidMethod(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/hello", nil)

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestForgotHandlerGet(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "Forgot"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestForgotHandlerPostMissingEmail(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/forgot", nil)

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_MISSING_EMAIL
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestForgotHandlerPostInvalidEmail(t *testing.T) {
	data := url.Values{"email": {"bad email"}}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/forgot", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_NO_SUCH_USER
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestForgotHandlerPostValidEmail(t *testing.T) {
	data := url.Values{"email": {"test@email"}}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/forgot", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "If you provided a valid email"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

/*
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
*/
