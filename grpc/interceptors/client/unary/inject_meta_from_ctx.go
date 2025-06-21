package unary

import (
	"context"

	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"

	"github.com/gw-gong/gwkit-go/grpc/interceptors/meta_data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func InjectMetaFromCtx() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := gwkit_common.GetRequestIDFromCtx(ctx)
		traceID := gwkit_common.GetTraceIDFromCtx(ctx)
		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {
			md = setMetaDataRequestID(md, requestID)
			md = setMetaDataTraceID(md, traceID)
		} else {
			md = metadata.New(map[string]string{})
			md = setMetaDataRequestID(md, requestID)
			md = setMetaDataTraceID(md, traceID)
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func setMetaDataRequestID(md metadata.MD, requestID string) metadata.MD {
	if requestID != "" {
		md.Set(meta_data.MetaKeyRequestID, requestID)
	}
	return md
}

func setMetaDataTraceID(md metadata.MD, traceID string) metadata.MD {
	if traceID != "" {
		md.Set(meta_data.MetaKeyTraceID, traceID)
	}
	return md
}
