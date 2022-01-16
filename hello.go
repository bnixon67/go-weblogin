package main

import (
	"log"
	"net/http"
)

// HelloPageData record
type HelloPageData struct {
	Message string
	User    User
}

// HelloHandler prints a simple hello message
func (app *App) HelloHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidMethod(w, r, []string{http.MethodGet}) {
		log.Println("invalid method", r.Method)
		return
	}

	// get sessionToken from cookie, if it exists
	var sessionToken string
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println("error getting cookie", err)
			w.WriteHeader(http.StatusInternalServerError)
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
	err = app.tmpls.ExecuteTemplate(w, "hello.html", HelloPageData{Message: "", User: currentUser})
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
