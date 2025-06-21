package common

import (
	"context"
)

type ContextKeyRequestID struct{}
type ContextKeyTraceID struct{}

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID{}, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	if value := ctx.Value(ContextKeyRequestID{}); value != nil {
		if requestID, ok := value.(string); ok {
			return requestID
		}
	}
	return ""
}

func SetTraceIDToCtx(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID{}, traceID)
}

func GetTraceIDFromCtx(ctx context.Context) string {
	if value := ctx.Value(ContextKeyTraceID{}); value != nil {
		if traceID, ok := value.(string); ok {
			return traceID
		}
	}
	return ""
}
