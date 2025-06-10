package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyLogger struct{}

func SetLoggerToCtx(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger{}, logger)
}

func GetLoggerFromCtx(ctx context.Context) *zap.Logger {
    if logger, ok := ctx.Value(ctxKeyLogger{}).(*zap.Logger); ok {
        return logger
    }
    return zap.L()
}