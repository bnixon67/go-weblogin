package main

import (
	"errors"
	"log"
	"net/http"
)

// HelloPageData contains data passed to the HTML template.
type HelloPageData struct {
	Message string
	User    User
}

// HelloHandler prints a simple hello message.
func (app *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
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
		currentUser, err = GetUserForSessionToken(app.db, sessionToken)
		if err != nil {
			log.Printf("failed to get user for session %q: %v", sessionToken, err)
			currentUser = User{}
			// delete invalid sessionToken to prevent session fixation
			http.SetCookie(w, &http.Cookie{Name: "sessionToken", Value: "", MaxAge: -1})
		} else {
			log.Println("UserName =", currentUser.UserName)
		}
	}

	// display page
	err = ExecTemplateOrError(app.tmpls, w, "hello.html", HelloPageData{Message: "", User: currentUser})
	if err != nil {
		log.Printf("error executing template: %v", err)
		return
	}
}
