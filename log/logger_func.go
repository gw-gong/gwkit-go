package log

import (
	"context"

	"go.uber.org/zap"
)

// ================================ Logger Functions ================================

func Debug(msg string, fields ...field) {
	zap.L().Debug(msg, fields...)
}

func Info(msg string, fields ...field) {
	zap.L().Info(msg, fields...)
}

func Warn(msg string, fields ...field) {
	zap.L().Warn(msg, fields...)
}

func Error(msg string, fields ...field) {
	zap.L().Error(msg, fields...)
}

func Debugc(ctx context.Context, msg string, fields ...field) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Debug(msg, fields...)
}

func Infoc(ctx context.Context, msg string, fields ...field) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Info(msg, fields...)
}

func Warnc(ctx context.Context, msg string, fields ...field) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Warn(msg, fields...)
}

func Errorc(ctx context.Context, msg string, fields ...field) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Error(msg, fields...)
}

// ================================ Sugar Functions ================================

func Debugf(format string, args ...interface{}) {
	zap.L().Sugar().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	zap.L().Sugar().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	zap.L().Sugar().Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	zap.L().Sugar().Errorf(format, args...)
}

func Debugfc(ctx context.Context, format string, args ...interface{}) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Sugar().Debugf(format, args...)
}

func Infofc(ctx context.Context, format string, args ...interface{}) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Sugar().Infof(format, args...)
}

func Warnfc(ctx context.Context, format string, args ...interface{}) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Sugar().Warnf(format, args...)
}

func Errorfc(ctx context.Context, format string, args ...interface{}) {
	loggerFromCtx := getLoggerFromCtx(ctx)
	loggerFromCtx.Sugar().Errorf(format, args...)
}

// ================================ With Fields Functions ================================

func WithFields(ctx context.Context, fields ...field) context.Context {
	loggerFromCtx := getLoggerFromCtx(ctx)
	return setLoggerToCtx(ctx, loggerFromCtx.With(fields...))
}

func WithFieldRequestID(ctx context.Context, requestID string) context.Context {
	loggerFromCtx := getLoggerFromCtx(ctx)
	return setLoggerToCtx(ctx, loggerFromCtx.With(Str("rid", requestID)))
}

func WithFieldTraceID(ctx context.Context, traceID string) context.Context {
	loggerFromCtx := getLoggerFromCtx(ctx)
	return setLoggerToCtx(ctx, loggerFromCtx.With(Str("trace_id", traceID)))
}