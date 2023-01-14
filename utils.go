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
package weblogin

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
)

// GenerateRandomString returns n bytes encoded in URL friendly base64.
func GenerateRandomString(n int) (string, error) {
	// buffer to store n bytes
	b := make([]byte, n)

	// get b random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// convert to URL friendly base64
	return base64.URLEncoding.EncodeToString(b), err
}

var ErrNoRequest = errors.New("request is nil")

// GetCookieValue returns the Value for the named cookie or empty string if not found or other error.
func GetCookieValue(r *http.Request, name string) (string, error) {
	var value string
	if r == nil {
		return value, ErrNoRequest
	}

	cookie, err := r.Cookie(name)
	if err != nil {
		// ignore ErrNoCookie
		if !errors.Is(err, http.ErrNoCookie) {
			return value, err
		}
	} else {
		value = cookie.Value
	}

	return value, nil
}

const SessionTokenCookieName = "session"

// GetUser returns the current User or empty User if the session is not found.
func GetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) (User, error) {
	var user User

	// get sessionToken from cookie, if it exists
	sessionToken, err := GetCookieValue(r, SessionTokenCookieName)
	if err != nil {
		return user, err
	}

	// get user if there is a sessionToken
	if sessionToken != "" {
		user, err = GetUserForSessionToken(db, sessionToken)
		if err != nil {
			// delete invalid token to prevent session fixation
			http.SetCookie(w,
				&http.Cookie{
					Name:   SessionTokenCookieName,
					Value:  "",
					MaxAge: -1,
				})
		}
		// ignore session not found errors
		if errors.Is(err, ErrSessionNotFound) {
			err = nil
		}
	}

	return user, err
}
