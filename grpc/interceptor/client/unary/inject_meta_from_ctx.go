package unary

import (
	"context"

	"github.com/gw-gong/gwkit-go/util/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func InjectMetaFromCtx() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := trace.GetRequestIDFromCtx(ctx)
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			md = setMetaDataRequestID(md, requestID)
		} else {
			md = metadata.New(map[string]string{})
			md = setMetaDataRequestID(md, requestID)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func setMetaDataRequestID(md metadata.MD, requestID string) metadata.MD {
	if requestID != "" {
		md.Set(trace.LoggerFieldRequestID, requestID)
	}
	return md
}
