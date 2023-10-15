package worker

import (
	"fmt"
	"github.com/charmbracelet/log"
)

const (
	taskProcessorLoggerName = "task processor"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Print(level log.Level, args ...interface{}) {
	switch level {
	case log.DebugLevel:
		log.Debug(taskProcessorLoggerName, "details", fmt.Sprint(args...))
	case log.InfoLevel:
		log.Info(taskProcessorLoggerName, "details", fmt.Sprint(args...))
	case log.WarnLevel:
		log.Warn(taskProcessorLoggerName, "details", fmt.Sprint(args...))
	case log.ErrorLevel:
		log.Error(taskProcessorLoggerName, "details", fmt.Sprint(args...))
	case log.FatalLevel:
		log.Fatal(taskProcessorLoggerName, "details", fmt.Sprint(args...))
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(log.DebugLevel, args...)

}

func (l *Logger) Info(args ...interface{}) {
	l.Print(log.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(log.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(log.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Print(log.FatalLevel, args...)
}
