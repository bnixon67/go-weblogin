package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

type LogWriter struct{}

const timeFormat = "2006-01-02 15:04:05 "

func (w LogWriter) Write(data []byte) (int, error) {
	return fmt.Print(time.Now().Format(timeFormat), funcName(4), ": ", string(data))
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

// logPanicIfEmpty calls log.Panic with the given message if str is empty
func logPanicIfEmpty(str, message string) {
	if str == "" {
		log.Panic(message)
	}
}

func logIfEmpty(str, message string) bool {
	if str == "" {
		log.Print(message)
		return true
	}

	return false
}
