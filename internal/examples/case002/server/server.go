package main

import (
	"context"

	pb "github.com/gw-gong/gwkit-go/internal/examples/case002/protobuf"
	"github.com/gw-gong/gwkit-go/log"
)

func NewTestService() pb.TestServiceServer {
	return &testService{}
}

type testService struct {
	pb.UnimplementedTestServiceServer
}

func (s *testService) TestFunc(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {
	log.Infoc(ctx, "TestFunc", log.Str("request", request.RequestName))
	return &pb.TestResponse{ResponseMsg: "test"}, nil
}
