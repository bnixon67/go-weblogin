package main

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user stored in the database
type User struct {
	UserName       string
	SessionToken   string
	FullName       string
	Email          string
	SessionExpires time.Time
}

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrNoUserForEmail      = errors.New("no username for email")
	ErrNoUserForResetToken = errors.New("no username for resetToken")
	ErrTooManyRows         = errors.New("too many rows affected")
)

// GetUserForSessionToken returns a user for the given sessionToken
func GetUserForSessionToken(db *sql.DB, sessionToken string) (User, error) {
	user := User{}

	qry := `SELECT users.userName, value, fullName, email, expires FROM users INNER JOIN tokens ON users.userName=tokens.userName WHERE tokens.type = "session" AND value=?`
	result := db.QueryRow(qry, sessionToken)
	err := result.Scan(&user.UserName, &user.SessionToken, &user.FullName, &user.Email, &user.SessionExpires)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrSessionNotFound
		}
		return User{}, err
	}

	return user, err
}

// RowExists return true if the given query returns at least one row
func RowExists(db *sql.DB, qry string, args ...interface{}) (bool, error) {
	var num int

	row := db.QueryRow(qry, args...)
	err := row.Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, err
}

// UserExists returns true if the given userName already exists in db
func UserExists(db *sql.DB, userName string) (bool, error) {
	return RowExists(db, "SELECT 1 FROM users WHERE userName=?", userName)
}

// EmailExists returns true if the given email already exists
func EmailExists(db *sql.DB, email string) (bool, error) {
	return RowExists(db, "SELECT 1 FROM users WHERE email=?", email)
}

// GetUserNameForEmail returns the userName for a given email
func GetUserNameForEmail(db *sql.DB, email string) (string, error) {
	var userName string

	row := db.QueryRow("SELECT username FROM users WHERE email=?", email)
	err := row.Scan(&userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNoUserForEmail
		}
		return "", err
	}

	return userName, err
}

// GetUserNameForResetToken returns the userName for a given reset token
func GetUserNameForResetToken(db *sql.DB, tokenValue string) (string, error) {
	var userName string

	qry := `SELECT userName FROM tokens WHERE type="reset" AND value=?`
	row := db.QueryRow(qry, tokenValue)
	err := row.Scan(&userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNoUserForResetToken
		}
		return "", err
	}

	return userName, err
}

// CompareUserPassword compares the password and hashed password for the user.
// Returns nil on success or an error on failure.
func CompareUserPassword(db *sql.DB, userName, password string) error {
	// get hashed password for the given user
	qry := `SELECT hashedPassword FROM users WHERE username=?`
	result := db.QueryRow(qry, userName)
	var hashedPassword string
	err := result.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoSuchUser
		}
		return err
	}

	// compared hashed password with given password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
