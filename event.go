package weblogin

import (
	"database/sql"
	"log/slog"
	"time"
)

const (
	EventLogin     = "login"
	EventLogout    = "logout"
	EventRegister  = "register"
	EventSaveToken = "save token"
)

type Event struct {
	Name     string // name of the event
	Result   bool   // result of the event
	UserName string // username for the event
	Message  string // message associated with event
	Created  time.Time
}

// WriteEvent will write an event to the database. There is no return value and if an error is encountered, it will be logged.
func WriteEvent(db *sql.DB, name string, result bool, user, msg string) {
	logger := slog.With(slog.Group("event",
		slog.String("Name", name),
		slog.Bool("Result", result),
		slog.String("Message", msg),
		slog.String("UserName", user),
	))

	qry := `INSERT INTO events(userName, action, result, message) VALUES(?, ?, ?, ?)`
	_, err := db.Exec(qry, user, name, result, msg)
	if err != nil {
		logger.Error("could not WriteEvent", "err", err)
	}
	logger.Debug("WriteEvent")
}
