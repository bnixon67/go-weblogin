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

	// get sessionTokenValue from cookie, if it exists
	var sessionTokenValue string
	c, err := r.Cookie("sessionToken")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println("error getting cookie", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		sessionTokenValue = c.Value
	}

	// create an empty sessionToken cookie with negative MaxAge to delete
	http.SetCookie(w, &http.Cookie{Name: "sessionToken", Value: "", MaxAge: -1})

	// remove session from database
	// TODO: consider removing all sessions for user
	if sessionTokenValue != "" {
		err := RemoveToken(app.db, "session", sessionTokenValue)
		if err != nil {
			log.Printf("remove token failed for %q: %v", sessionTokenValue, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	}

	// display page
	err = app.tmpls.ExecuteTemplate(w, "logout.html", nil)
	if err != nil {
		log.Println("error executing template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
