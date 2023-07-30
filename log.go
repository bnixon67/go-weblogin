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
package weblogin

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
)

// LogWriter is a custom io.Writer to output log entries prefixed with date/time and the function name of the caller to log.
type LogWriter struct {
	w io.Writer
}

const (
	timeFormat = "2006-01-02 15:04:05 "
	fileMode   = 0o600
)

// InitLog initializes logging for the application.
func InitLog(logFileName string) error {
	var err error

	// configure log writter
	var w io.Writer
	if logFileName == "" {
		w = os.Stderr
	} else {
		w, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("InitLog: failed to open log file: %w", err)
		}
	}

	// configure logger
	// TODO: use config value for options
	opts := &slog.HandlerOptions{
		AddSource: true,
		// ReplaceAttr: redact,
	}
	logger := slog.New(slog.NewJSONHandler(w, opts))
	slog.SetDefault(logger)

	return nil
}

func redact(_ []string, a slog.Attr) slog.Attr {
	app, ok := a.Value.Any().(*App)
	if ok {
		app.Config.SMTPPassword = "[REDACTED]"
		return slog.Any(a.Key, app)
	}

	config, ok := a.Value.Any().(Config)
	if ok {
		config.SMTPPassword = "[REDACTED]"
		return slog.Any(a.Key, config)
	}

	return a
}
