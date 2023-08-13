/*
Copyright 2023 Bill Nixon

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
	"log/slog"
	"os"
)

const fileMode = 0o600

// InitLog initializes logging for the application.
func InitLog(logFileName string, logLevel slog.Level, addSource bool) error {
	var err error

	// configure log writter
	var w io.Writer
	if logFileName == "" {
		w = os.Stderr
	} else {
		w, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
		if err != nil {
			return fmt.Errorf("InitLog: failed to open log file: %w", err)
		}
	}

	// configure logger
	opts := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(w, opts))
	slog.SetDefault(logger)

	return nil
}
