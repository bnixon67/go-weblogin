package main

import (
	"log"
	"net/http"
	"time"
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

	// check for valid db
	if app.db == nil {
		log.Println("app.db is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check for valid tmpls
	if app.tmpls == nil {
		log.Println("app.tmpls is nil")
		w.WriteHeader(http.StatusInternalServerError)
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
		} else {
			// check if token is expired
			if currentUser.SessionExpires.Before(time.Now()) {
				log.Printf("token expired for %q", currentUser.UserName)
				currentUser = User{}
			} else {
				log.Println("UserName =", currentUser.UserName)
			}
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
