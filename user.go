/*
Copyright 2022 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user stored in the database.
type User struct {
	UserName        string
	FullName        string
	Email           string
	Admin           bool
	Created         time.Time
	LastLoginTime   time.Time
	LastLoginResult string
}

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrNoUserForEmail      = errors.New("no username for email")
	ErrNoUserForResetToken = errors.New("no username for resetToken")
	ErrSessionExpired      = errors.New("session expired")
	ErrGetLastLoginFailed  = errors.New("failed to get last login")
)

// GetUserForSessionToken returns a user for the given sessionToken.
func GetUserForSessionToken(db *sql.DB, sessionToken string) (User, error) {
	var (
		expires time.Time
		user    User
	)

	hashedValue := hash(sessionToken)

	qry := `SELECT users.userName, fullName, email, expires, admin FROM users INNER JOIN tokens ON users.userName=tokens.userName WHERE tokens.type = "session" AND hashedValue=? LIMIT 1`
	result := db.QueryRow(qry, hashedValue)
	err := result.Scan(&user.UserName, &user.FullName, &user.Email, &expires, &user.Admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrSessionNotFound
		}
		return User{}, err
	}

	if expires.Before(time.Now()) {
		return User{}, ErrSessionExpired
	}

	user.LastLoginTime, user.LastLoginResult, err = LastLoginForUser(db, user.UserName)
	if err != nil {
		return user, fmt.Errorf("%w: %v", ErrGetLastLoginFailed, err)
	}

	return user, err
}

// GetUserForName returns a user for the given userName.
func GetUserForName(db *sql.DB, userName string) (User, error) {
	var user User

	qry := `SELECT userName, fullName, email, admin FROM users WHERE userName=? LIMIT 1`
	result := db.QueryRow(qry, userName)
	err := result.Scan(&user.UserName, &user.FullName, &user.Email, &user.Admin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, err
}

// RowExists return true if the given query returns at least one row.
func RowExists(db *sql.DB, qry string, args ...interface{}) (bool, error) {
	var num int

	row := db.QueryRow(qry, args...)
	err := row.Scan(&num)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, err
}

// UserExists returns true if the given userName already exists in db.
func UserExists(db *sql.DB, userName string) (bool, error) {
	return RowExists(db, "SELECT 1 FROM users WHERE userName=? LIMIT 1", userName)
}

// EmailExists returns true if the given email already exists.
func EmailExists(db *sql.DB, email string) (bool, error) {
	return RowExists(db, "SELECT 1 FROM users WHERE email=? LIMIT 1", email)
}

// GetUserNameForEmail returns the userName for a given email.
func GetUserNameForEmail(db *sql.DB, email string) (string, error) {
	var userName string

	row := db.QueryRow("SELECT username FROM users WHERE email=?", email)
	err := row.Scan(&userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoUserForEmail
		}
		return "", err
	}

	return userName, err
}

// GetUserNameForResetToken returns the userName for a given reset token.
func GetUserNameForResetToken(db *sql.DB, tokenValue string) (string, error) {
	var userName string
	hashedValue := hash(tokenValue)

	qry := `SELECT userName FROM tokens WHERE type="reset" AND hashedValue=?`
	row := db.QueryRow(qry, hashedValue)
	err := row.Scan(&userName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoUserForResetToken
		}
		return "", err
	}

	return userName, err
}

var ErrNoSuchUser = errors.New("no such user")

// CompareUserPassword compares the password and hashed password for the user.
// Returns nil on success or an error on failure.
func CompareUserPassword(db *sql.DB, userName, password string) error {
	// get hashed password for the given user
	qry := `SELECT hashedPassword FROM users WHERE username=? LIMIT 1`
	result := db.QueryRow(qry, userName)

	var hashedPassword string
	err := result.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

// RegisterUser registers a user with the given values.
// Returns nil on success or an error on failure.
func RegisterUser(db *sql.DB, userName, fullName, email, password string) error {
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// store the user and hashed password
	_, err = db.Exec("INSERT INTO users(username, hashedPassword, fullName, email) VALUES (?, ?, ?, ?)",
		userName, hashedPassword, fullName, email)
	if err != nil {
		return err
	}

	return nil
}

func LastLoginForUser(db *sql.DB, userName string) (time.Time, string, error) {
	var lastLogin time.Time
	var result string

	// get the second row, if it exists, since first row is current login
	qry := `SELECT created, result FROM events WHERE userName = ? AND action = ? ORDER BY created DESC LIMIT 1,1`

	row := db.QueryRow(qry, userName, ActionLogin)
	err := row.Scan(&lastLogin, &result)
	if err != nil {
		// ignore ErrNoRows since there may not be a last login
		if errors.Is(err, sql.ErrNoRows) {
			return lastLogin, result, nil
		}
		return lastLogin, result, err
	}

	return lastLogin, result, err
}
