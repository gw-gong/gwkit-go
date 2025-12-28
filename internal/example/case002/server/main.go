package main

import (
	"context"
	"fmt"
	"net"

	"github.com/gw-gong/gwkit-go/grpc/consul"
	"github.com/gw-gong/gwkit-go/grpc/interceptor/server/unary"
	"github.com/gw-gong/gwkit-go/internal/example/case002/protobuf"
	"github.com/gw-gong/gwkit-go/log"
	"github.com/gw-gong/gwkit-go/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	ServiceName = "test_service"
	ServiceTag  = "test"

	ServerPort = 8080
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	util.ExitOnErr(context.Background(), err)
	defer syncFn()

	consulClient, err := consul.NewConsulClient(consul.DefaultConsulAgentAddr)
	util.ExitOnErr(context.Background(), err)

	registerEntry := &consul.RegisterEntry{
		ServiceName: ServiceName,
		Tags:        []string{ServiceTag},
	}
	registerEntry.GenerateServiceID()
	err = consulClient.Register(registerEntry, ServerPort, false)
	util.ExitOnErr(context.Background(), err)
	defer func() {
		err = consulClient.Deregister(registerEntry.ServiceID)
		if err != nil {
			log.Error("consul service deregistration failed", log.Err(err))
		}
		log.Info("consul service deregistration succeeded")
	}()

	log.Info("consul service registration succeeded")

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			unary.PanicRecoverInterceptor(),
			unary.ParseMetaToCtx(),
		),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	testService := NewTestService()
	protobuf.RegisterTestServiceServer(grpcServer, testService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ServerPort))
	util.ExitOnErr(context.Background(), err)
	defer listener.Close()

	log.Info("服务启动成功", log.Str("port", fmt.Sprintf("%d", ServerPort)))
	err = grpcServer.Serve(listener)
	util.ExitOnErr(context.Background(), err)
}
