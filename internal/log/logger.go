package log

import "log"

type LoggingFunc func(format string, v ...interface{})

var printf LoggingFunc = log.Printf
var fatalf LoggingFunc = log.Fatalf

func SetPrintf(f LoggingFunc) LoggingFunc {
	prev := printf
	printf = f

	return prev
}

func SetFatalf(f LoggingFunc) LoggingFunc {
	prev := fatalf
	fatalf = f

	return prev
}

func Printf(format string, v ...interface{}) {
	printf(format, v...)
}
