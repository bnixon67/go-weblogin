package main

import (
	"testing"
)

func TestConfigIsValid(t *testing.T) {
	cases := []struct {
		config   Config
		expected bool
	}{
		{Config{SQLDriverName: "", SQLDataSourceName: "", ParseGlobPattern: ""}, false},
		{Config{SQLDriverName: "", SQLDataSourceName: "", ParseGlobPattern: "foo"}, false},
		{Config{SQLDriverName: "", SQLDataSourceName: "foo", ParseGlobPattern: ""}, false},
		{Config{SQLDriverName: "", SQLDataSourceName: "foo", ParseGlobPattern: "foo"}, false},
		{Config{SQLDriverName: "foo", SQLDataSourceName: "", ParseGlobPattern: ""}, false},
		{Config{SQLDriverName: "foo", SQLDataSourceName: "", ParseGlobPattern: "foo"}, false},
		{Config{SQLDriverName: "foo", SQLDataSourceName: "foo", ParseGlobPattern: ""}, false},
		{Config{SQLDriverName: "foo", SQLDataSourceName: "foo", ParseGlobPattern: "foo"}, true},
	}

	for _, testCase := range cases {
		t.Logf("%+v", testCase.config)
		got := testCase.config.IsValid()
		if got != testCase.expected {
			t.Errorf("c.IsValid(%+v) = %v; expected %v", testCase.config, got, testCase.expected)
		}
	}
}

func TestNewConfigFromFile(t *testing.T) {
	// test empty (invaild) file name
	_, err := NewConfigFromFile("")
	if err == nil {
		t.Errorf("NewConfigFromFile for empty filename is nil")
	}

	var fileName string

	// test with a valid filename and file with empty json
	fileName = "testdata/empty.json"
	config, err := NewConfigFromFile(fileName)
	if err != nil {
		t.Fatalf("NewConfigFromFile(%q) failed: %v", fileName, err)
	}
	if config != (Config{}) {
		t.Errorf("got %+v, expected %+v", app.config, Config{})
	}

	// test with a valid filename and file with invalid json
	fileName = "testdata/invalid.json"
	config, err = NewConfigFromFile(fileName)
	if err == nil {
		t.Errorf("expected NewConfigFromFile(%q) to fail", fileName)
	}
	if config != (Config{}) {
		t.Errorf("got %+v, expected %+v", config, Config{})
	}

	// test with a valid filename and file with valid json
	fileName = "testdata/valid.json"
	config, err = NewConfigFromFile(fileName)
	if err != nil {
		t.Fatalf("NewConfigFromFile(%q) failed: %v", fileName, err)
	}
	expectedConfig := Config{
		SQLDriverName:     "testSQLDriverName",
		SQLDataSourceName: "testSQLDataSourceName",
		ParseGlobPattern:  "testParseGlobPattern",
	}
	if config != expectedConfig {
		t.Errorf("got %+v, expected %+v", config, expectedConfig)
	}
}
