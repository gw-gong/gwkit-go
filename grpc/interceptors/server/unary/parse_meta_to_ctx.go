package unary

import (
	"context"

	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/utils/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ParseMetaToCtx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			requestIDs := md.Get(trace.LoggerFieldRequestID)
			if len(requestIDs) > 0 {
				requestID := requestIDs[0]
				if requestID == "" {
					requestID = trace.GenerateRequestID()
				}
				ctx = trace.SetRequestIDToCtx(ctx, requestID)
				ctx = log.WithFieldRequestID(ctx, requestID)
			}
		}
		return handler(ctx, req)
	}
}
