package weblogin

import (
	"database/sql"
	"log/slog"
	"time"
)

const (
	ActionLogin     = "login"
	ActionLogout    = "logout"
	ActionRegister  = "register"
	ActionSaveToken = "save token"
)

type Event struct {
	UserName string
	Created  time.Time
	Action   string
	Result   bool
	Message  string
}

// WriteEvent will write an event to the database. There is no return value and if an error is encountered, it will be logged.
func WriteEvent(db *sql.DB, event Event) {
	qry := `INSERT INTO events(userName, action, result, message) VALUES(?, ?, ?, ?)`
	_, err := db.Exec(qry, event.UserName, event.Action, event.Result, event.Message)
	if err != nil {
		slog.Error("could not WriteEvent", "event", event, "err", err)
	}
}
