package common

import "context"

type ContextKeyRequestID struct{}
type ContextKeyTraceID struct{}

func SetRequestIDToCtx(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID{}, requestID)
}

func GetRequestIDFromCtx(ctx context.Context) string {
	return ctx.Value(ContextKeyRequestID{}).(string)
}

func SetTraceIDToCtx(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ContextKeyTraceID{}, traceID)
}

func GetTraceIDFromCtx(ctx context.Context) string {
	return ctx.Value(ContextKeyTraceID{}).(string)
}
