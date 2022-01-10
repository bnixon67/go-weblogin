package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRegisterHandlerInvalidMethod(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/register", nil)

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestRegisterHandlerGet(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/register", nil)

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "Register"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestRegisterHandlerPostMissingValues(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", nil)

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_REGISTER_MISSING_VALUES
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestRegisterHandlerPostExistingUser(t *testing.T) {
	data := url.Values{
		"userName":  {"test"},
		"fullName":  {"full name"},
		"email":     {"email"},
		"password1": {"password one"},
		"password2": {"password two"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_REGISTER_USER_EXISTS
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestRegisterHandlerPostExistingEmail(t *testing.T) {
	randomUserName, err := GenerateRandomString(8)
	if err != nil {
		t.Fatalf("could not GenerateRandomString")
	}
	data := url.Values{
		"userName":  {randomUserName},
		"fullName":  {"full name"},
		"email":     {"test@email"},
		"password1": {"password one"},
		"password2": {"password two"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_REGISTER_EMAIL_EXISTS
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestRegisterHandlerPostMismatchedPassword(t *testing.T) {
	randomUserName, err := GenerateRandomString(8)
	if err != nil {
		t.Fatalf("could not GenerateRandomString")
	}
	data := url.Values{
		"userName":  {randomUserName},
		"fullName":  {"full name"},
		"email":     {randomUserName + "@email"},
		"password1": {"password one"},
		"password2": {"password two"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MSG_REGISTER_MISMATCHED_PASSWORDS
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestRegisterHandlerPostValid(t *testing.T) {
	randomUserName, err := GenerateRandomString(8)
	if err != nil {
		t.Fatalf("could not GenerateRandomString")
	}
	data := url.Values{
		"userName":  {randomUserName},
		"fullName":  {"full name"},
		"email":     {randomUserName + "@email"},
		"password1": {"password"},
		"password2": {"password"},
	}

	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(data.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.RegisterHandler(w, r)

	expectedStatus := http.StatusSeeOther
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := ""
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}

	got := w.Header().Get("Location")
	expected := "/login"
	if got != expected {
		t.Fatalf("got location %q, expected %q", got, expected)
	}
}
