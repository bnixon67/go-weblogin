/*
Copyright 2022 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin_test

import (
	"errors"
	"math"
	"reflect"
	"testing"

	weblogin "github.com/bnixon67/go-weblogin"
	"github.com/google/go-cmp/cmp"
)

func TestNewConfigFromFile(t *testing.T) {
	testCases := []struct {
		name           string
		configFileName string
		wantErr        error
		wantConfig     weblogin.Config
	}{
		{
			name:           "emptyFileName",
			configFileName: "",
			wantErr:        weblogin.ErrConfigOpen,
			wantConfig:     weblogin.Config{},
		},
		{
			name:           "emptyJSON",
			configFileName: "testdata/empty.json",
			wantErr:        nil,
			wantConfig:     weblogin.Config{},
		},
		{
			name:           "invalidJSON",
			configFileName: "testdata/invalid.json",
			wantErr:        weblogin.ErrConfigDecode,
			wantConfig:     weblogin.Config{},
		},
		{
			name:           "validJSON",
			configFileName: "testdata/valid.json",
			wantErr:        nil,
			wantConfig: weblogin.Config{
				Title:               "Test Title",
				ServerHost:          "test host",
				ServerPort:          "test port",
				BaseURL:             "test URL",
				SQLDriverName:       "testSQLDriverName",
				SQLDataSourceName:   "testSQLDataSourceName",
				ParseGlobPattern:    "testParseGlobPattern",
				SessionExpiresHours: 42,
				SMTPHost:            "test SMTP host",
				SMTPPort:            "test SMTP port",
				SMTPUser:            "test SMTP user",
				SMTPPassword:        "test SMTP password",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(*testing.T) {
			config, err := weblogin.NewConfigFromFile(tc.configFileName)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got err %q, want %q for NewConfigFromFile(%q)", err, tc.wantErr, tc.configFileName)
			}

			if diff := cmp.Diff(config, tc.wantConfig); diff != "" {
				t.Errorf("config did not match (-got +want):\n%s", diff)
			}
		})
	}
}

func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func TestConfigIsValid(t *testing.T) {
	type tcase struct {
		config   weblogin.Config
		expected bool
	}

	var cases []tcase

	// required fields
	required := []string{
		"Title",
		"ServerHost",
		"ServerPort",
		"BaseURL",
		"SQLDriverName",
		"SQLDataSourceName",
		"ParseGlobPattern",
		"SMTPHost",
		"SMTPPort",
		"SMTPUser",
		"SMTPPassword",
	}

	// generate test cases based on required fields by looping thru all the possibilities and using bit logic to set fields
	for a := 0; a < int(math.Pow(2, float64(len(required)))); a++ {
		config := weblogin.Config{}

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

func TestConfigMarshalJSON(t *testing.T) {
	testCases := []struct {
		name  string
		input weblogin.Config
		want  string
	}{
		{
			name: "test",
			input: weblogin.Config{
				Title:             "AppConfig",
				SQLDataSourceName: "user:password@localhost/db",
				SMTPPassword:      "supersecret",
			},
			want: `{"Title":"AppConfig","ServerHost":"","ServerPort":"","BaseURL":"","SQLDriverName":"","SQLDataSourceName":"[REDACTED]","ParseGlobPattern":"","SessionExpiresHours":0,"SMTPHost":"","SMTPPort":"","SMTPUser":"","SMTPPassword":"[REDACTED]"}`,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.MarshalJSON()
			if err != nil {
				t.Fatalf("Error during MarshalJSON: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got\n%s\n, want\n%s\n", got, tc.want)
			}
		})
	}
}

func TestConfigString(t *testing.T) {
	testCases := []struct {
		name  string
		input weblogin.Config
		want  string
	}{
		{
			name: "test",
			input: weblogin.Config{
				Title:             "AppConfig",
				SQLDataSourceName: "user:password@localhost/db",
				SMTPPassword:      "supersecret",
			},
			want: `{Title:AppConfig ServerHost: ServerPort: BaseURL: SQLDriverName: SQLDataSourceName:[REDACTED] ParseGlobPattern: SessionExpiresHours:0 SMTPHost: SMTPPort: SMTPUser: SMTPPassword:[REDACTED]}`,
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.String()
			if got != tc.want {
				t.Errorf("got\n%s\n, want\n%s\n", got, tc.want)
			}
		})
	}
}
