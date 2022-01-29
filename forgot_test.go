/*
   Copyright 2022 Bill Nixon

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
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
		t.Errorf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}
}

func TestForgotHandlerGet(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Errorf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "Forgot"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Errorf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}

func TestForgotHandlerPostMissingEmail(t *testing.T) {
	app := AppForTest(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/forgot", nil)

	app.ForgotHandler(w, r)

	expectedStatus := http.StatusOK
	if w.Code != expectedStatus {
		t.Errorf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MsgMissingEmail
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Errorf("got body %q, expected %q in body", w.Body, expectedInBody)
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
		t.Errorf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := MsgNoSuchUser
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Errorf("got body %q, expected %q in body", w.Body, expectedInBody)
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
		t.Errorf("got status %d %q, expected %d %q", w.Code, http.StatusText(w.Code), expectedStatus, http.StatusText(expectedStatus))
	}

	expectedInBody := "If you provided a valid email"
	if !strings.Contains(w.Body.String(), expectedInBody) {
		t.Errorf("got body %q, expected %q in body", w.Body, expectedInBody)
	}
}
