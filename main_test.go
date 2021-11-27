package main

import (
	"os"
	"testing"
)

var app *App

func TestMain(m *testing.M) {
	// TODO: either define a default config.json or pass in as parameter
	var err error

	app, err = NewApp("config.json")
	if err != nil {
		return
	}

	os.Exit(m.Run())
}
