package main

import (
	"testing"
)

var app *App

// AppForTest is a helper function that returns an App used for testing.
func AppForTest(t *testing.T) *App {
	var err error

	if app == nil {
		app, err = NewApp("config.json", "test.log")
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
			configFileName: "config.json",
			logFileName:    "test.log",
			errExpected:    false,
			appExpected:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(*testing.T) {
			app, err := NewApp(c.configFileName, c.logFileName)
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
