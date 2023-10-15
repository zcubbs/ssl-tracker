package worker

import (
	"github.com/zcubbs/logwrapper/logger"
)

const (
	taskProcessorLoggerName = "task processor"
)

type Logger struct {
	logger logger.Logger
}

func NewLogger(logger logger.Logger) *Logger {
	return &Logger{logger}
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(taskProcessorLoggerName, args...)

}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(taskProcessorLoggerName, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(taskProcessorLoggerName, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(taskProcessorLoggerName, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(taskProcessorLoggerName, args...)
}
