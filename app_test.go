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
package weblogin_test

import (
	"errors"
	"testing"

	weblogin "github.com/bnixon67/go-weblogin"
	_ "github.com/go-sql-driver/mysql"
)

// global to provide a singleton app.
var app *weblogin.App //nolint

const TestLogFile = "test.log"

// AppForTest is a helper function that returns an App used for testing.
func AppForTest(t *testing.T) *weblogin.App {
	if app == nil {
		var err error
		app, err = weblogin.NewApp("testdata/test_config.json", TestLogFile)
		if err != nil {
			app = nil

			t.Fatalf("cannot create NewApp, %v", err)
		}
	}

	return app
}

// TestNewApp provides tests for the NewApp function.
func TestNewApp(t *testing.T) {
	testCases := []struct {
		name           string
		configFileName string
		logFileName    string
		wantErr        error
		isAppExpected  bool
	}{
		{
			name:           "validConfigAndLog",
			configFileName: "testdata/test_config.json",
			logFileName:    "test.log",
			wantErr:        nil,
			isAppExpected:  true,
		},
		{
			name:           "emptyConfigFileName",
			configFileName: "",
			logFileName:    "test.log",
			wantErr:        weblogin.ErrOpenConfig,
			isAppExpected:  false,
		},
		{
			name:           "badLogFileName",
			configFileName: "",
			logFileName:    "/foo/bar",
			wantErr:        weblogin.ErrInitLog,
			isAppExpected:  false,
		},
		{
			name:           "emptyConfig",
			configFileName: "testdata/empty.json",
			logFileName:    "test.log",
			wantErr:        weblogin.ErrInvalidConfig,
			isAppExpected:  false,
		},
		{
			name:           "invalidDB",
			configFileName: "testdata/invalid_db.json",
			logFileName:    "test.log",
			wantErr:        weblogin.ErrInitDB,
			isAppExpected:  false,
		},
		{
			name:           "invalidTemplates",
			configFileName: "testdata/invalid_tmpl.json",
			logFileName:    "test.log",
			wantErr:        weblogin.ErrInitTemplates,
			isAppExpected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(*testing.T) {
			app, err := weblogin.NewApp(tc.configFileName, tc.logFileName)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got err %q, want %q for NewApp(%q, %q)", err, tc.wantErr, tc.configFileName, tc.logFileName)
			}

			gotApp := app != nil
			if gotApp != tc.isAppExpected {
				t.Errorf("gotApp is %t, want %t for NewApp(%q, %q)", gotApp, tc.isAppExpected, tc.configFileName, tc.logFileName)
			}
		})
	}
}
