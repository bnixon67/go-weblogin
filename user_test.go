package weblogin_test

import (
	"errors"
	"testing"
	"time"

	weblogin "github.com/bnixon67/go-weblogin"
)

func TestLastLoginForUser(t *testing.T) {
	var zeroTime time.Time

	dt := time.Date(2023, time.January, 15, 1, 0, 0, 0, time.UTC)

	cases := []struct {
		userName string
		want     time.Time
		err      error
	}{
		{"no such user", zeroTime, nil},
		{"test1", zeroTime, nil},
		{"test2", dt, nil},
		{"test3", dt.Add(time.Hour), nil},
		{"test4", dt.Add(time.Hour * 2), nil},
	}

	app := AppForTest(t)

	for _, tc := range cases {
		got, _, err := weblogin.LastLoginForUser(app.DB, tc.userName)
		if !errors.Is(err, tc.err) {
			t.Errorf("LastLoginForUser(db, %q)\ngot err '%v' want '%v'", tc.userName, err, tc.err)
		}
		if got != tc.want {
			t.Errorf("LastLoginForUser(db, %q)\n got '%v'\nwant '%v'", tc.userName, got, tc.want)
		}
	}
}
