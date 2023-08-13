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
	EventSaveToken = "save_token"
	EventReset     = "reset_pass"
	EventMax       = "1234567890"
)

type Event struct {
	Name     string // name of the event
	Result   bool   // result of the event
	UserName string // username for the event
	Message  string // message associated with event
	Created  time.Time
}

// WriteEvent will write an event to the database. There is no return value and if an error is encountered, it will be logged.
func WriteEvent(db *sql.DB, name string, result bool, user, message string) {
	logger := slog.With(slog.Group("event",
		slog.String("Name", name),
		slog.Bool("Result", result),
		slog.String("Message", message),
		slog.String("UserName", user),
	))

	qry := `INSERT INTO events(name, result, userName, message) VALUES(?, ?, ?, ?)`
	_, err := db.Exec(qry, name, result, user, message)
	if err != nil {
		logger.Error("could not WriteEvent", "err", err)
	}
	logger.Debug("WriteEvent")
}
