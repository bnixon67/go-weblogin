package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogWriter struct{}

const timeFormat = "2006-01-02 15:04:05 "

func (LogWriter) Write(data []byte) (int, error) {
	return fmt.Fprint(os.Stderr, time.Now().Format(timeFormat), funcName(4), ": ", string(data))
}

func funcName(depth int) string {
	pc, _, _, ok := runtime.Caller(depth)

	if !ok {
		return "runtime.Caller(1) not ok"
	}

	fn := runtime.FuncForPC(pc)

	if fn == nil {
		return "?()"
	}

	names := strings.Split(fn.Name(), ".")
	return names[len(names)-1]
}

// logIfEmpty logs message if str is empty and returns true, otherwise false
func logIfEmpty(str, message string) bool {
	if str == "" {
		log.Print(message)
		return true
	}

	return false
}
