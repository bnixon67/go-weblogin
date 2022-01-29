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
	"log"
	"os"
	"runtime"
	"strings"
	"time"
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
	// use custom writer for log
	lw, err := NewLogWriter(logFileName)
	if err != nil {
		return fmt.Errorf("InitLog: %w", err)
	}
	log.SetFlags(0)
	log.SetOutput(lw)

	return err
}

// NewLogWriter creates a new LogWriter. The filename defines where to write the logfile. If filename is blank, then os.Stderr is used.
func NewLogWriter(filename string) (LogWriter, error) {
	var lw LogWriter
	var err error

	if filename == "" {
		lw.w = os.Stderr
	} else {
		lw.w, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fileMode)
	}

	return lw, err
}

// Write satisfies io.Writer interface to output with prefixed as noted above.
func (lw LogWriter) Write(data []byte) (int, error) {
	return fmt.Fprint(lw.w, time.Now().Format(timeFormat), FuncName(4), ": ", string(data))
}

// FuncName returns the function name at the depth provided.
func FuncName(depth int) string {
	// get program counter
	pc, _, _, ok := runtime.Caller(depth)

	if !ok {
		return "runtime.Caller() failed"
	}

	// get function for program counter
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?()"
	}

	// fn.Name() may inclue the package, so return only the function name
	names := strings.Split(fn.Name(), ".")
	return names[len(names)-1]
}
