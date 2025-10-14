package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/grpc/interceptors/client/unary"
	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
	gwkit_str "github.com/gw-gong/gwkit-go/utils/str"
	"github.com/gw-gong/gwkit-go/utils/trace"

	"google.golang.org/grpc"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	gwkit_common.ExitOnErr(context.Background(), err)
	defer syncFn()

	requestID := gwkit_str.GenerateULID()
	ctx := trace.SetRequestIDToCtx(context.Background(), requestID)
	ctx = log.WithFieldRequestID(ctx, requestID)

	testClient, err := NewTestClient("127.0.0.1:8500", "test_service", "test", "",
		grpc.WithChainUnaryInterceptor(
			unary.InjectMetaFromCtx(),
		),
	)
	gwkit_common.ExitOnErr(ctx, err)

	_, _ = testClient.TestFunc(ctx, "test")

	log.Info("rpc 调用完成")
}
