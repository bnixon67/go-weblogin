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

// LogWriter is a custom io.Writer that can be used to output log entries. Each write will be prefixed with date/time and the function name of the caller to log.
type LogWriter struct {
	w io.Writer
}

const timeFormat = "2006-01-02 15:04:05 "

// NewLogWriter creates a new LogWriter. The filename defines where to write the logfile. If filename is blank, then os.Stderr is used.
func NewLogWriter(filename string) (LogWriter, error) {
	var lw LogWriter
	var err error

	if filename == "" {
		lw.w = os.Stderr
	} else {
		lw.w, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	return lw, err
}

// Write satifies the io.Writer interface and outputs the data prefixed as noted above.
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

// logIfEmpty logs message if str is empty and returns true, otherwise false.
func logIfEmpty(str, message string) bool {
	if str == "" {
		log.Print(message)
		return true
	}

	return false
}
