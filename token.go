package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"
)

// Token represent a token for the user
type Token struct {
	Value   string
	Expires time.Time
	Type    string
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// SaveNewToken creates and saves a token for user of size that expires in hrs
func SaveNewToken(db *sql.DB, tType, userName string, size, hrs int) (Token, error) {
	var err error

	token := Token{Type: tType}
	token.Value, err = GenerateRandomString(size)
	if err != nil {
		return Token{}, err
	}
	token.Expires = time.Now().Add(time.Duration(hrs) * time.Hour)

	// hash the token to avoid reuse if database is compromised
	hashedValue := hash(token.Value)

	qry := `INSERT INTO tokens(hashedValue, expires, type, userName) VALUES(?, ?, ?, ?)`
	_, err = db.Exec(qry, hashedValue, token.Expires, tType, userName)
	return token, err
}

// RemoveToken removes the given sessionToken
func RemoveToken(db *sql.DB, tType, tValue string) error {
	hashedValue := hash(tValue)

	qry := `DELETE FROM tokens WHERE type = ? AND hashedValue = ?`
	_, err := db.Exec(qry, tType, hashedValue)
	return err
}
