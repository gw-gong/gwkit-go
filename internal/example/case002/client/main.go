package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/grpc/consul"
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
	ctx = trace.WithLogFieldRequestID(ctx, requestID)

	consulClient, err := consul.NewConsulClient(consul.DefaultConsulAgentAddr)
	util.ExitOnErr(ctx, err)

	testClient, err := NewTestClient(consulClient, &consul.HealthyGrpcConnEntry{
		ServiceName: "test_service",
		Tag:         "test",
		Opts:        []grpc.DialOption{grpc.WithChainUnaryInterceptor(unary.InjectMetaFromCtx())},
	})
	util.ExitOnErr(ctx, err)

	_, _ = testClient.TestFunc(ctx, "test")

	log.Info("rpc 调用完成")
}
