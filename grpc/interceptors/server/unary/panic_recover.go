package unary

import (
	"context"

	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"

	"google.golang.org/grpc"
)

func PanicRecoverInterceptor() grpc.ServerOption {
	unaryServerInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		gwkit_common.WithRecover(func() {
			resp, err = handler(ctx, req)
		}, gwkit_common.WithPanicHandler(func(err interface{}) {
			gwkit_common.DefaultPanicWithCtx(ctx, err)
		}))
		return
	}
	return grpc.UnaryInterceptor(unaryServerInterceptor)
}
