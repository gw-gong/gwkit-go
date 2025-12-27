package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/grpc/consul"
	"github.com/gw-gong/gwkit-go/internal/example/case002/protobuf"
	"github.com/gw-gong/gwkit-go/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient interface {
	TestFunc(ctx context.Context, requestName string) (responseMsg string, err error)
}

func NewTestClient(agentAddr, serviceName, tag string, opts ...grpc.DialOption) (TestClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // 使用 insecure 连接 (不使用 TLS, 开发环境使用)
	return newTestClient(agentAddr, serviceName, tag, opts...)
}

type testClient struct {
	client protobuf.TestServiceClient
}

func newTestClient(agentAddr, serviceName, tag string, opts ...grpc.DialOption) (TestClient, error) {
	conn, err := consul.NewHealthyGrpcConn(agentAddr, serviceName, tag, opts...)
	if err != nil {
		return nil, err
	}
	return &testClient{client: protobuf.NewTestServiceClient(conn)}, nil
}

func (c *testClient) TestFunc(ctx context.Context, requestName string) (responseMsg string, err error) {
	request := &protobuf.TestRequest{RequestName: requestName}
	response, err := c.client.TestFunc(ctx, request)
	if err != nil {
		log.Errorc(ctx, "TestFunc", log.Err(err))
		return "", err
	}
	log.Infoc(ctx, "TestFunc", log.Str("response", response.ResponseMsg))
	return response.ResponseMsg, nil
}
