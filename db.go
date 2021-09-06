package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// initDB initializes a connection to the database
func initDB(driverName, dataSourceName string) (*sql.DB, error) {
	log.Print("Initialize database connection")

	// open connection to database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// set desire connection parameters
	// TODO: move values to config file
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Ping to confirm connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

type User struct {
	UserName       string
	SessionToken   string
	FirstName      string
	LastName       string
	Email          string
	SessionExpires time.Time
}

// GetUserForSessionToken returns a user for the given sessionToken
func GetUserForSessionToken(sessionToken string) (User, error) {
	user := User{}

	result := db.QueryRow(
		"SELECT userName, sessionToken, firstName, lastName, email, sessionExpires FROM users WHERE sessionToken=?",
		sessionToken)
	err := result.Scan(
		&user.UserName, &user.SessionToken, &user.FirstName, &user.LastName, &user.Email, &user.SessionExpires)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No result for sessionToken %q", sessionToken)
			return user, errors.New("Invalid sessionToken")
		}
		log.Println("Query for sessionToken failed", err)
		return user, errors.New("Query failed")
	}

	return user, err
}

// CheckForUserName returns true if the given userName already exists
func CheckForUserName(userName string) (bool, error) {
	var name string

	result := db.QueryRow("SELECT userName FROM users WHERE userName=?", userName)
	err := result.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Println("Query for userName failed", err)
		return false, err
	}

	return true, err
}
