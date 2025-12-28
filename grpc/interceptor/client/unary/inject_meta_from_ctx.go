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
		traceID := trace.GetTraceIDFromCtx(ctx)
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			md = setMetaDataTraceInfo(md, requestID, traceID)
		} else {
			md = metadata.New(map[string]string{})
			md = setMetaDataTraceInfo(md, requestID, traceID)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func setMetaDataTraceInfo(md metadata.MD, requestID string, traceID string) metadata.MD {
	if requestID != "" {
		md.Set(trace.LoggerFieldRequestID, requestID)
	}
	if traceID != "" {
		md.Set(trace.LoggerFieldTraceID, traceID)
	}
	return md
}
