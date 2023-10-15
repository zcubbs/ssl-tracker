package logger

import (
	"github.com/zcubbs/logwrapper"
	"github.com/zcubbs/logwrapper/logger"
)

const loggerName = "tlz"

var baseLogger logger.Logger

func init() {
	baseLogger = New()
}

func L() logger.Logger {
	return baseLogger
}

func New() logger.Logger {
	return logwrapper.NewLogger(
		logger.CharmLoggerType,
		loggerName,
		logger.TextFormat,
	)
}

func SetFormat(format string) {
	baseLogger.SetFormat(format)
}

func SetLevel(level string) {
	baseLogger.SetLevel(level)
}
