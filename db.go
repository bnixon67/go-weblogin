package main

import (
	"database/sql"
	"log"
	"time"
)

// InitDB initializes a connection to the database.
func InitDB(driverName, dataSourceName string) (*sql.DB, error) {
	log.Println("initialize database connection")

	// open connection to database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// set desire connection parameters
	// TODO: move values to config file
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// ping database to confirm connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
