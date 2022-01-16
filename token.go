package main

import (
	"database/sql"
	"time"
)

// Token represent a token for the user
type Token struct {
	Value   string
	Expires time.Time
	Type    string
}

// SaveNewToken creates and saves a new session for user that expires in hrs
func SaveNewToken(db *sql.DB, tType, userName string, hrs int) (Token, error) {
	var err error

	token := Token{Type: tType}
	token.Value, err = GenerateRandomString(32)
	if err != nil {
		return Token{}, err
	}
	token.Expires = time.Now().Add(time.Duration(hrs) * time.Hour)

	qry := "INSERT INTO tokens(value, expires, type, userName) VALUES(?, ?, ?, ?)"
	_, err = db.Exec(qry, token.Value, token.Expires, tType, userName)
	return token, err
}

// RemoveToken removes the given sessionToken
func RemoveToken(db *sql.DB, tType, tValue string) error {
	qry := "DELETE FROM tokens WHERE type = ? AND value = ?"
	_, err := db.Exec(qry, tType, tValue)
	return err
}
