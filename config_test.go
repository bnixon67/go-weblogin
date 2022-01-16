package main

import (
	"math"
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	InitLogging("test.log")

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
		t.Errorf("NewConfigFromFile(%q) failed: %v", fileName, err)
	}
	if config != (Config{}) {
		t.Errorf("got %+v, expected %+v", config, Config{})
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
		t.Errorf("NewConfigFromFile(%q) failed: %v", fileName, err)
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

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func TestConfigIsValid(t *testing.T) {
	InitLogging("test.log")

	type tcase struct {
		config   Config
		expected bool
	}
	var cases []tcase

	// required fields
	required := []string{"ServerHost", "ServerPort", "SQLDriverName", "SQLDataSourceName", "ParseGlobPattern"}

	// generate test cases based on required fields by looping thru all the possibilities and using bit logic to set fields
	for a := 0; a < int(math.Pow(2, float64(len(required)))); a++ {
		config := Config{}
		for n := len(required) - 1; n >= 0; n-- {
			if hasBit(a, uint(n)) {
				reflect.ValueOf(&config).Elem().FieldByName(required[n]).SetString("x")
			}
		}
		cases = append(cases, tcase{config, false})
	}
	// last case should be true since all required fields are present
	cases[len(cases)-1].expected = true

	for _, testCase := range cases {
		got, _ := testCase.config.IsValid()
		if got != testCase.expected {
			t.Errorf("c.IsValid(%+v) = %v; expected %v", testCase.config, got, testCase.expected)
		}
	}
}
