package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// LogoutHandler handles /logout requests
func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
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

	// create an empty sessionToken cookie with negative MaxAge to delete
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionToken",
		Value:  "",
		MaxAge: -1,
	})

	// remove session from database
	// TODO: consider removing all sessions for user
	if sessionToken != "" {
		err := RemoveSession(app.db, sessionToken)
		if err != nil {
			log.Printf("remove session failed for %q: %v", sessionToken, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	// display page
	err = app.tmpls.ExecuteTemplate(w, "logout.html", nil)
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
