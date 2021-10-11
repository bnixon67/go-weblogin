package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// LogoutHandler handles /logout requests
func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "from", r.RemoteAddr)

	// only GET method is allowed
	if r.Method != "GET" {
		log.Println("Invalid method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// create an empty sessionToken cookie with negative MaxAge to delete
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionToken",
		Value:  "",
		MaxAge: -1,
	})

	// display page
	err := app.tmpls.ExecuteTemplate(w, "logout.html", nil)
	if err != nil {
		log.Println("error executing template", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
