package logger

import "fmt"

func Info(message string, args ...interface{}) {
	fmt.Printf("INFO  "+message+"\n", args...)
}

func Warn(message string, args ...interface{}) {
	fmt.Printf("WARN  "+message+"\n", args...)
}

func Error(message string, args ...interface{}) {
	fmt.Printf("ERROR "+message+"\n", args...)
}

func Success(message string, args ...interface{}) {
	fmt.Printf("OK    "+message+"\n", args...)
}