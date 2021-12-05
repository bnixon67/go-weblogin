package main

import (
	"testing"
)

func TestInitDB(t *testing.T) {
	// test invalid
	db, err := initDB("", "")
	if err == nil {
		t.Errorf("initDB returned nil err for empty values")
	}
	if db != nil {
		t.Errorf("initDB returned non-nil db for empty values")
	}

	// TODO: test valid
}
