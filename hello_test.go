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

	HelloHandler(w, r)

	expectedStatus := http.StatusMethodNotAllowed
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestHelloHandlerWithoutCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)

	HelloHandler(w, r)

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

	HelloHandler(w, r)

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
	// TODO: generate valid session token instead of hard coding
	token := "47MCgM9wrkkkfWZiWfFqBo5AI87u_NN4qPjjArH_osY="
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	r.AddCookie(&http.Cookie{Name: "sessionToken", Value: token})

	HelloHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Fatalf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "You must <a href=\"/login\">login</a>"
	if strings.Contains(w.Body.String(), expectedInBody) {
		t.Fatalf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}
