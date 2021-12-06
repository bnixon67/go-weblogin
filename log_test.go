package main

import (
	"testing"
)

func TestInitLogging(t *testing.T) {
	// invalid file name
	err := InitLogging("/foo/bar")
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

func TestLogIfEmpty(t *testing.T) {
	v := ""
	expected := true
	isEmpty := logIfEmpty(v, "test")
	if isEmpty != expected {
		t.Errorf("got %v, expected %v, for %q", isEmpty, expected, v)
	}

	v = "foo"
	expected = false
	isEmpty = logIfEmpty(v, "test")
	if isEmpty != expected {
		t.Errorf("got %v, expected %v, for %q", isEmpty, expected, v)
	}
}
