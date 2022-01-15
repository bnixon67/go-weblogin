package main

import (
	"database/sql"
	"time"
)

// Session represent a session for the user
type Session struct {
	Token   string
	Expires time.Time
}

// SaveNewSession creates and saves a new session for user that expires in hrs
func SaveNewSession(db *sql.DB, userName string, hrs int) (Session, error) {
	var err error

	session := Session{}
	session.Token, err = GenerateRandomString(32)
	if err != nil {
		return Session{}, err
	}
	session.Expires = time.Now().Add(time.Duration(hrs) * time.Hour)

	qry := "INSERT INTO sessions(token, expires, userName) VALUES(?, ?, ?)"
	_, err = db.Exec(qry, session.Token, session.Expires, userName)
	return session, err
}

// RemoveSession removes the given sessionToken
func RemoveSession(db *sql.DB, sessionToken string) error {
	qry := "DELETE FROM sessions WHERE token = ?"
	_, err := db.Exec(qry, sessionToken)
	return err
}
