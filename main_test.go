package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// TODO: either define a default config.json or pass in as parameter
	InitApp("config.json")
	os.Exit(m.Run())
}
