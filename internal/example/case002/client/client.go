package main

import (
	"context"
	"fmt"

	"github.com/gw-gong/gwkit-go/grpc/consul"
	"github.com/gw-gong/gwkit-go/internal/example/case002/protobuf"
	"github.com/gw-gong/gwkit-go/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TestClient interface {
	TestFunc(ctx context.Context, requestName string) (responseMsg string, err error)
}

type testClient struct {
	client protobuf.TestServiceClient
}

func NewTestClient(consulClient consul.ConsulClient, entry *consul.HealthyGrpcConnEntry) (TestClient, error) {
	if entry == nil {
		return nil, fmt.Errorf("entry is nil")
	}
	entry.Opts = append(entry.Opts, grpc.WithTransportCredentials(insecure.NewCredentials())) // intranet env does not use tls

	conn, err := consulClient.GetHealthyGrpcConn(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
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
