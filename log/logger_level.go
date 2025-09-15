package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LoggerLevelDebug string = "debug"
	LoggerLevelInfo  string = "info"
	LoggerLevelWarn  string = "warn"
	LoggerLevelError string = "error"
)

func MapLoggerLevel(level string) zapcore.Level {
	switch level {
	case LoggerLevelDebug:
		return zap.DebugLevel
	case LoggerLevelInfo:
		return zap.InfoLevel
	case LoggerLevelWarn:
		return zap.WarnLevel
	case LoggerLevelError:
		return zap.ErrorLevel
	}
	Errorf("unknown logger level: %s", level)
	return zap.InfoLevel
}
