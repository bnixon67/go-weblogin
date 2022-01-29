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
package weblogin_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	weblogin "github.com/bnixon67/go-web-login"
)

func TestStringContains(t *testing.T) {
	type test struct {
		arr    []string
		val    string
		expect bool
	}

	tests := []test{
		{arr: []string{}, val: "", expect: false},
		{arr: []string{http.MethodGet}, val: http.MethodGet, expect: true},
		{arr: []string{http.MethodGet, http.MethodPost}, val: http.MethodGet, expect: true},
		{arr: []string{http.MethodGet}, val: http.MethodPost, expect: false},
		{arr: []string{http.MethodGet, http.MethodPatch}, val: http.MethodPost, expect: false},
	}

	for _, tc := range tests {
		got := weblogin.StringContains(tc.arr, tc.val)
		if got != tc.expect {
			t.Errorf("got %v, expected %v, for StringContains(%q, %q)", got, tc.expect, tc.arr, tc.val)
		}
	}
}

func TestCheckMethods(t *testing.T) {
	type test struct {
		arr    []string
		val    string
		expect bool
		status int
		inBody string
		allow  string
	}

	tests := []test{
		{
			arr:    []string{http.MethodGet},
			val:    http.MethodGet,
			expect: true,
			status: http.StatusOK,
			inBody: "",
			allow:  "",
		},
		{
			arr:    []string{http.MethodGet, http.MethodPost},
			val:    http.MethodGet,
			expect: true,
			status: http.StatusOK,
			inBody: "",
			allow:  "",
		},
		{
			arr:    []string{http.MethodGet},
			val:    http.MethodPost,
			expect: false,
			status: http.StatusMethodNotAllowed,
			inBody: http.StatusText(http.StatusMethodNotAllowed),
			allow:  strings.Join([]string{http.MethodGet, http.MethodOptions}, ", "),
		},
		{
			arr:    []string{http.MethodGet, http.MethodPatch},
			val:    http.MethodPost,
			expect: false,
			status: http.StatusMethodNotAllowed,
			inBody: http.StatusText(http.StatusMethodNotAllowed),
			allow:  strings.Join([]string{http.MethodGet, http.MethodPatch, http.MethodOptions}, ", "),
		},
		{
			arr:    []string{http.MethodGet, http.MethodPatch},
			val:    http.MethodOptions,
			expect: false,
			status: http.StatusNoContent,
			inBody: "",
			allow:  strings.Join([]string{http.MethodGet, http.MethodPatch, http.MethodOptions}, ", "),
		},
	}

	for _, tc := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(tc.val, "/", nil)

		got := weblogin.ValidMethod(w, r, tc.arr)
		body := w.Body.String()
		allow := w.Header().Get("Allow")
		if got != tc.expect || w.Code != tc.status || !strings.Contains(body, tc.inBody) || w.Header().Get("Allow") != tc.allow {
			t.Errorf("got %v %v %q %q, expected %v %v %q %q, for %q, %q", got, w.Code, body, allow, tc.expect, tc.status, tc.inBody, tc.allow, tc.arr, tc.val)
		}
	}
}
