package unary

import (
	"context"

	"github.com/gw-gong/gwkit-go/grpc/interceptors/meta_data"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ParseMetaToCtx() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			requestID := md.Get(meta_data.MetaKeyRequestID)
			if len(requestID) > 0 {
				requestID := requestID[0]
				ctx = gwkit_common.SetRequestIDToCtx(ctx, requestID)
				ctx = log.WithFieldRequestID(ctx, requestID)
			}
			traceID := md.Get(meta_data.MetaKeyTraceID)
			if len(traceID) > 0 {
				traceID := traceID[0]
				ctx = gwkit_common.SetTraceIDToCtx(ctx, traceID)
				ctx = log.WithFieldTraceID(ctx, traceID)
			}
		}
		return handler(ctx, req)
	}
}
