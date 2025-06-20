package main

import (
	"context"

	pb "github.com/gw-gong/gwkit-go/examples/case002/internal/protobuf"
	consul_agent "github.com/gw-gong/gwkit-go/grpc/consul/agent"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient interface {
	TestFunc(ctx context.Context, requestName string) (responseMsg string, err error)
}

func NewTestClient(agentAddr, serviceName, tag string, dc string, opts ...grpc.DialOption) (TestClient, error) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // 使用 insecure 连接 (不使用 TLS, 开发环境使用)
	return newTestClient(agentAddr, serviceName, tag, dc, opts...)
}

type testClient struct {
	client pb.TestServiceClient
}

func newTestClient(agentAddr, serviceName, tag string, dc string, opts ...grpc.DialOption) (TestClient, error) {
	conn, err := consul_agent.NewHealthyGrpcConn(agentAddr, serviceName, tag, dc, opts...)
	if err != nil {
		return nil, err
	}
	return &testClient{client: pb.NewTestServiceClient(conn)}, nil
}

func (c *testClient) TestFunc(ctx context.Context, requestName string) (responseMsg string, err error) {
	request := &pb.TestRequest{RequestName: requestName}
	response, err := c.client.TestFunc(ctx, request)
	if err != nil {
		return "", err
	}
	return response.ResponseMsg, nil
}
