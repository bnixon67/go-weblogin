package main

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

// InitLogging initializes logging for the application.
func InitLogging(logFileName string) error {
	// use custom writer for log
	lw, err := NewLogWriter(logFileName)
	if err != nil {
		log.Printf("unable to create NewLogWriter, %v", err)
		return err
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
	return fmt.Fprint(lw.w, time.Now().Format(timeFormat), funcName(4), ": ", string(data))
}

// funcName returns the function name at the depth provided.
func funcName(depth int) string {
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
