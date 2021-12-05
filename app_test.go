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
			t.Fatal("cannot create NewApp")
		}
	}

	return app
}

func TestNewApp(t *testing.T) {
	app, err := NewApp("", "test.log")
	if err == nil {
		t.Error("expected non-nill err for NewApp(\"\", \"\")\n")
	}
	if app != nil {
		t.Errorf("got app=%v, expected nil for NewApp(\"\", \"\")\n", app)
	}
}
