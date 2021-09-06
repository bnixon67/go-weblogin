package main

import (
	"log"
	"runtime"
)

func funcName() string {
	pc, _, _, ok := runtime.Caller(1)

	if !ok {
		return "runtime.Caller(1) not ok"
	}

	fn := runtime.FuncForPC(pc)

	if fn == nil {
		return "runtime.FuncForPC nil"
	} else {
		return fn.Name()
	}
}

// logPanicIsEmpty calls log.Panic with the given message if str is empty
func logPanicIsEmpty(str, message string) {
	if str == "" {
		log.Panic(message)
	}
}
