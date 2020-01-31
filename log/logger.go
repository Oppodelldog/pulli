package log

import "log"

type LogFunc func(format string, v ...interface{})

var printf LogFunc
var fatalf LogFunc

func init() {
	printf = log.Printf
	fatalf = log.Fatalf
}

func SetPrintf(f LogFunc) LogFunc {
	prev := printf
	printf = f

	return prev
}

func SetFatalf(f LogFunc) LogFunc {
	prev := fatalf
	fatalf = f

	return prev
}

func Printf(format string, v ...interface{}) {
	printf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	fatalf(format, v...)
}
