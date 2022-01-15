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

	qry := "SELECT userName, sessionToken, fullName, email, sessionExpires FROM users WHERE sessionToken=?"
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

// CheckForUserName returns true if the given userName already exists
func (app *App) CheckForUserName(userName string) (bool, error) {
	var num int

	row := app.db.QueryRow("SELECT 1 FROM users WHERE userName=?", userName)
	err := row.Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, err
}

// CheckForEmail returns true if the given email already exists
func (app *App) CheckForEmail(email string) (bool, error) {
	var num int

	row := app.db.QueryRow("SELECT 1 FROM users WHERE email=?", email)
	err := row.Scan(&num)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, err
}

// GetUserNameForEmail returns the userName for a given email
func (app *App) GetUserNameForEmail(email string) (string, error) {
	var userName string

	row := app.db.QueryRow("SELECT username FROM users WHERE email=?", email)
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
func (app *App) GetUserNameForResetToken(resetToken string) (string, error) {
	var userName string

	row := app.db.QueryRow("SELECT username FROM users WHERE resetToken=?", resetToken)
	err := row.Scan(&userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNoUserForResetToken
		}
		return "", err
	}

	return userName, err
}

func (app *App) SaveResetTokenForUser(userName, resetToken string) error {
	result, err := app.db.Exec("UPDATE users SET resetToken  = ? WHERE username = ?", resetToken, userName)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrTooManyRows
	}

	return err
}

func (app *App) CheckUserPassword(userName, password string) error {
	// get hashed password for the given user
	result := app.db.QueryRow("SELECT hashedPassword FROM users WHERE username=?", userName)
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
