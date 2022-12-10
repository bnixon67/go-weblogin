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
	"database/sql"
	"errors"
	"log"
	"net/http"
)

// UsersPageData contains data passed to the HTML template.
type UsersPageData struct {
	Title   string
	Message string
	User    User
	Users   []User
}

// UsersHandler prints a simple hello message.
func (app *App) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet}) {
		log.Println("invalid method", r.Method)
		return
	}

	// get sessionToken from cookie, if it exists
	var sessionToken string
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			log.Println("error getting cookie", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		sessionToken = c.Value
	}

	// get user for sessionToken
	var currentUser User
	if sessionToken != "" {
		currentUser, err = GetUserForSessionToken(app.DB, sessionToken)
		if err != nil {
			log.Printf("failed to get user for session %q: %v", sessionToken, err)
			currentUser = User{}
			// delete invalid sessionToken to prevent session fixation
			http.SetCookie(w, &http.Cookie{Name: "sessionToken", Value: "", MaxAge: -1})
		} else {
			log.Println("UserName =", currentUser.UserName)
		}
	}

	users, err := GetUsers(app.DB)
	if err != nil {
		log.Printf("failed to GetUsers: %v", err)
	}

	// display page
	err = RenderTemplate(app.Tmpls, w, "users.html",
		UsersPageData{Message: "", User: currentUser, Users: users})
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}

// GetUsers returns a list of all users.
func GetUsers(db *sql.DB) ([]User, error) {
	var users []User
	var err error

	if db == nil {
		log.Print("db is nil")
		return users, errors.New("invdalid db")
	}

	qry := `SELECT userName, fullName, email, admin, created FROM users`

	rows, err := db.Query(qry)
	if err != nil {
		log.Printf("query for users failed, %v", err)
		return users, err

	}
	defer rows.Close()

	for rows.Next() {
		var user User

		err = rows.Scan(&user.UserName, &user.FullName, &user.Email, &user.Admin, &user.Created)
		if err != nil {
			log.Printf("rows.Scan failed, %v", err)
		}

		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("rows.Err failed, %v", err)
	}

	return users, err
}
