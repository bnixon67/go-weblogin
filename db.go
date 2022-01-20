package main

import (
	"database/sql"
	"fmt"
	"time"
)

// InitDB initializes a connection to the database.
func InitDB(driverName, dataSourceName string) (*sql.DB, error) {
	// open connection to database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("InitDB: %w", err)
	}

	// set desire connection parameters
	// TODO: move values to config file
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// ping database to confirm connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("InitDB: %w", err)
	}

	return db, err
}
