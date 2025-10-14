package trace

import (
	"context"
)

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
