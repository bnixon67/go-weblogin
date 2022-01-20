package main

import (
	"testing"
)

const TestLogFile = "test.log"

func TestInitLogging(t *testing.T) {
	// invalid file name
	err := InitLog("/foo/bar")
	if err == nil {
		t.Errorf("got nil, expected non-nil for InitLogging with invalid file name")
	}
}

func TestFuncName(t *testing.T) {
	expected := "TestFuncName"
	name := funcName(1)
	if name != expected {
		t.Errorf("got %q, expected %q", name, expected)
	}
}
