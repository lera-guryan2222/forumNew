package logger

import "log"

type Logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type stdLogger struct{}

func NewLogger() Logger {
	return &stdLogger{}
}

func (l *stdLogger) Info(args ...interface{}) {
	log.Println("INFO:", args)
}

func (l *stdLogger) Error(args ...interface{}) {
	log.Println("ERROR:", args)
}
