package unary

import (
	"context"

	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicRecoverInterceptor() grpc.ServerOption {
	unaryServerInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, errInterceptor error) {
		gwkit_common.WithRecover(func() {
			resp, errInterceptor = handler(ctx, req)
		}, gwkit_common.WithPanicHandler(func(err interface{}) {
			gwkit_common.DefaultPanicWithCtx(ctx, err)
			errInterceptor = status.Errorf(codes.Internal, "Internal Server Error (panic recovered)")
			resp = nil
		}))
		return
	}
	return grpc.UnaryInterceptor(unaryServerInterceptor)
}
