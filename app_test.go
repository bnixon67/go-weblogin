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
	"testing"

	weblogin "github.com/bnixon67/go-weblogin"
	_ "github.com/go-sql-driver/mysql"
)

// global to provide a singleton app.
var app *weblogin.App //nolint

// AppForTest is a helper function that returns an App used for testing.
func AppForTest(t *testing.T) *weblogin.App {
	if app == nil {
		var err error
		app, err = weblogin.NewApp("test_config.json", TestLogFile)
		if err != nil {
			app = nil

			t.Fatalf("cannot create NewApp, %v", err)
		}
	}

	return app
}

// TestNewApp provides tests for the NewApp function.
func TestNewApp(t *testing.T) {
	cases := []struct {
		name           string
		configFileName string
		logFileName    string
		errExpected    bool
		appExpected    bool
	}{
		{
			name:           "emptyConfigFileName",
			configFileName: "",
			logFileName:    "test.log",
			errExpected:    true,
			appExpected:    false,
		},
		{
			name:           "badLogFileName",
			configFileName: "",
			logFileName:    "/foo/bar",
			errExpected:    true,
			appExpected:    false,
		},
		{
			name:           "emptyConfig",
			configFileName: "testdata/empty.json",
			logFileName:    "test.log",
			errExpected:    true,
			appExpected:    false,
		},
		{
			name:           "validConfigAndLog",
			configFileName: "test_config.json",
			logFileName:    "test.log",
			errExpected:    false,
			appExpected:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(*testing.T) {
			app, err := weblogin.NewApp(c.configFileName, c.logFileName)
			if c.errExpected && err == nil {
				t.Errorf("expected error, got err==nil for NewApp(%q, %q)", c.configFileName, c.logFileName)
			}
			if !c.errExpected && err != nil {
				t.Errorf("expected no error, got err=%q for NewApp(%q, %q)", err, c.configFileName, c.logFileName)
			}
			if c.appExpected && app == nil {
				t.Errorf("expected app, got app=nil for NewApp(%q, %q)", c.configFileName, c.logFileName)
			}
			if !c.appExpected && app != nil {
				t.Errorf("expected app==nil, got app=%v for NewApp(%q, %q)", app, c.configFileName, c.logFileName)
			}
		})
	}
}
