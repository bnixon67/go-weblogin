package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
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
	ErrSessionTokenNotFound = errors.New("sessionToken not found")
	ErrNoUserForEmail       = errors.New("no username for email")
	ErrNoUserForResetToken  = errors.New("no username for resetToken")
	ErrTooManyRows          = errors.New("too many rows affected")
)

// GetUserForSessionToken returns a user for the given sessionToken
func (app *App) GetUserForSessionToken(sessionToken string) (User, error) {
	user := User{}

	qry := "SELECT userName, sessionToken, fullName, email, sessionExpires FROM users WHERE sessionToken=?"
	result := app.db.QueryRow(qry, sessionToken)
	err := result.Scan(&user.UserName, &user.SessionToken, &user.FullName, &user.Email, &user.SessionExpires)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrSessionTokenNotFound
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
		log.Printf("query for userName %q failed: %v", userName, err)
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
		log.Printf("query for email %q failed: %v", email, err)
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
			log.Printf("no username for %q", email)
			return "", ErrNoUserForEmail
		}
		log.Printf("query for email %q failed: %v", email, err)
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
			log.Printf("No username for %q", resetToken)
			return "", ErrNoUserForResetToken
		}
		log.Printf("query for resetToken %q failed: %v", resetToken, err)
		return "", err
	}

	return userName, err
}

func (app *App) SaveResetTokenForUser(userName, resetToken string) error {
	result, err := app.db.Exec("UPDATE users SET resetToken  = ? WHERE username = ?", resetToken, userName)
	if err != nil {
		log.Printf("update resetToken for %q failed: %v", userName, err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected() is not nil: %v", err)
		return err
	}
	if rowsAffected != 1 {
		log.Printf("expected to affect 1 row, affected %d", rowsAffected)
		return ErrTooManyRows
	}

	log.Printf("Saved resetToken %q for %q", resetToken, userName)
	return err
}
