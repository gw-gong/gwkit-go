package unary

import (
	"context"

	"github.com/gw-gong/gwkit-go/util/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ParseMetaToCtx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			var havaTraceInfo bool
			requestIDs := md.Get(trace.LoggerFieldRequestID)
			if len(requestIDs) > 0 && requestIDs[0] != "" {
				ctx = trace.SetRequestIDToCtx(ctx, requestIDs[0])
				ctx = trace.WithLogFieldRequestID(ctx, requestIDs[0])
				havaTraceInfo = true
			}
			traceIDs := md.Get(trace.LoggerFieldTraceID)
			if len(traceIDs) > 0 && traceIDs[0] != "" {
				ctx = trace.SetTraceIDToCtx(ctx, traceIDs[0])
				ctx = trace.WithLogFieldTraceID(ctx, traceIDs[0])
				havaTraceInfo = true
			}
			if !havaTraceInfo {
				newRequestID := trace.GenerateRequestID()
				ctx = trace.SetRequestIDToCtx(ctx, newRequestID)
				ctx = trace.WithLogFieldRequestID(ctx, newRequestID)
			}
		}
		return handler(ctx, req)
	}
}
