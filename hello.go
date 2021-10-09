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
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// only GET method is allowed
	if r.Method != "GET" {
		log.Println("invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// check for valid db
	if db == nil {
		log.Println("db is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check for valid tmpls
	if tmpls == nil {
		log.Println("tmpls is nil")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get sessionToken from cookie, if it exists
	var sessionToken string
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("no sessionToken cookie")
		} else {

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
		currentUser, err = GetUserForSessionToken(sessionToken)
		if err != nil {
			log.Println("GetUserForSessionToken failed:", err)
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
	err = tmpls.ExecuteTemplate(w, "hello.html", HelloPageData{Message: "", User: currentUser})
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
