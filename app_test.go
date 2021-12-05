package main

import (
	"testing"
)

var app *App

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
