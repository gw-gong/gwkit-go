package main

import (
	"context"
	"fmt"
	"net"

	"github.com/gw-gong/gwkit-go/grpc/interceptor/server/unary"
	"github.com/gw-gong/gwkit-go/internal/example/case002/protobuf"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util/common"
	"github.com/gw-gong/gwkit-go/util/str"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	common.ExitOnErr(context.Background(), err)
	defer syncFn()

	serviceRegistry, err := NewTestServiceRegistry(ServiceName)
	common.ExitOnErr(context.Background(), err)
	serviceID := str.GenerateUUID()
	err = serviceRegistry.Register(serviceID, ServerPort, []string{ServiceTag})
	common.ExitOnErr(context.Background(), err)
	defer func() {
		err = serviceRegistry.Deregister(serviceID)
		if err != nil {
			log.Error("consul 服务注销失败", log.Err(err))
		}
		log.Info("consul 服务注销成功")
	}()

	log.Info("consul 服务注册成功")

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unary.PanicRecoverInterceptor(),
			unary.ParseMetaToCtx(),
		),
	)

	// 创建并注册健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// 注册服务
	testService := NewTestService()
	protobuf.RegisterTestServiceServer(grpcServer, testService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	common.ExitOnErr(context.Background(), err)
	defer listener.Close()

	log.Info("服务启动成功", log.Str("port", fmt.Sprintf("%d", ServerPort)))
	err = grpcServer.Serve(listener)
	common.ExitOnErr(context.Background(), err)
}
