package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/grpc/interceptor/client/unary"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util"
	"github.com/gw-gong/gwkit-go/util/str"
	"github.com/gw-gong/gwkit-go/util/trace"

	"google.golang.org/grpc"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	util.ExitOnErr(context.Background(), err)
	defer syncFn()

	requestID := str.GenerateULID()
	ctx := trace.SetRequestIDToCtx(context.Background(), requestID)
	ctx = log.WithFieldRequestID(ctx, requestID)

	testClient, err := NewTestClient("127.0.0.1:8500", "test_service", "test",
		grpc.WithChainUnaryInterceptor(
			unary.InjectMetaFromCtx(),
		),
	)
	util.ExitOnErr(ctx, err)

	_, _ = testClient.TestFunc(ctx, "test")

	log.Info("rpc 调用完成")
}
