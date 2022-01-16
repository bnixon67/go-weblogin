package main

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	// TODO: test valid
	// test invalid
	db, err := InitDB("", "")
	if err == nil {
		t.Errorf("initDB returned nil err for empty values")
	}
	if db != nil {
		t.Errorf("initDB returned non-nil db for empty values")
	}
}
