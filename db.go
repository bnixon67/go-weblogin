/*
Copyright 2023 Bill Nixon

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
)

var (
	ErrDBOpen = errors.New("failed to open db")
	ErrDBPing = errors.New("failed to ping db")
)

// InitDB initializes a connection to the database.
func InitDB(driverName, dataSourceName string) (*sql.DB, error) {
	fn := "InitDB"

	// open connection to database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", fn, ErrDBOpen, err)
	}

	// set desire connection parameters
	// TODO: move values to config file
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// ping database to confirm connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %v", fn, ErrDBPing, err)
	}

	return db, err
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
