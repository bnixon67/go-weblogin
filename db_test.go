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
package weblogin_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	weblogin "github.com/bnixon67/go-weblogin"
)

/*
func TestInitDB(t *testing.T) {
	// TODO: test valid
	// test invalid
	db, err := weblogin.InitDB("", "")
	if err == nil {
		t.Errorf("initDB returned nil err for empty values")
	}
	if db != nil {
		t.Errorf("initDB returned non-nil db for empty values")
	}
}
*/

// CustomMockDriver is a custom driver that only returns a connection error.
type CustomMockDriver struct{}

func (d CustomMockDriver) Open(name string) (driver.Conn, error) {
	if name == "valid_source" {
		return nil, nil
	}
	return nil, errors.New("mock connection error")
}

func TestInitDB(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		driverName     string
		dataSourceName string
		wantErr        error
	}{
		{
			name:           "valid",
			driverName:     "mock_driver",
			dataSourceName: "valid_source",
			wantErr:        nil,
		},
		{
			name:           "invalid driver",
			driverName:     "invalid_driver",
			dataSourceName: "valid_source",
			wantErr:        weblogin.ErrDBOpen,
		},
		{
			name:           "invalid source",
			driverName:     "mock_driver",
			dataSourceName: "invalid_source",
			wantErr:        weblogin.ErrDBPing,
		},
	}

	// Create a map to mock the sql.Open function
	mockDriver := &CustomMockDriver{}
	sql.Register("mock_driver", mockDriver)

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := weblogin.InitDB(tc.driverName, tc.dataSourceName)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got err %q, want %q for InitDB(%q, %q)", err, tc.wantErr, tc.driverName, tc.dataSourceName)
			}
		})
	}
}
