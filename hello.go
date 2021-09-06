package main

import (
	"log"
	"net/http"
)

// HelloHandler prints a simple hello message
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.Method)

	// only GET method is allowed
	if r.Method != "GET" {
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get sessionToken from cookie, if it exists
	var sessionToken string
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Print("No sessionToken cookie")
		} else {
			log.Println("Error getting cookie", err)
		}
	} else {
		sessionToken = c.Value
	}

	// get user for sessionToken
	var currentUser User
	if sessionToken != "" {
		currentUser, err = GetUserForSessionToken(sessionToken)
		if err != nil {
			log.Println("GetUserForSessionToken failed", err)
			return
		}
		log.Printf("%+v", currentUser)
	}

	// display page
	err = tmpls.ExecuteTemplate(w, "hello.html", LoginPageData{"", currentUser})
	if err != nil {
		log.Println("Error in executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
