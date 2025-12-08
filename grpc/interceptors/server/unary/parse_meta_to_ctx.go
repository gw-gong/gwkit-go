package unary

import (
	"context"

	"github.com/gw-gong/gwkit-go/log"
	gwkit_trace "github.com/gw-gong/gwkit-go/util/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ParseMetaToCtx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			requestIDs := md.Get(gwkit_trace.LoggerFieldRequestID)
			if len(requestIDs) > 0 {
				requestID := requestIDs[0]
				if requestID == "" {
					requestID = gwkit_trace.GenerateRequestID()
				}
				ctx = gwkit_trace.SetRequestIDToCtx(ctx, requestID)
				ctx = log.WithFieldRequestID(ctx, requestID)
			}
		}
		return handler(ctx, req)
	}
}
