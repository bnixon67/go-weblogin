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

	_, err = db.Query("UPDATE users SET sessionToken = ?, sessionExpires = ? WHERE username = ?", session.Token, session.Expires, userName)
	return session, err
}
