package logger

import (
	"github.com/zcubbs/logwrapper"
	"github.com/zcubbs/logwrapper/logger"
)

const loggerName = "tlz"

func GetLogger() logger.Logger {
	return GetLoggerWithName(loggerName)
}

func GetLoggerWithName(name string) logger.Logger {
	l := logwrapper.NewLogger(
		logger.CharmLoggerType,
		name,
		logger.TextFormat,
	)

	return l
}
