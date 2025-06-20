package main

import (
	"context"

	pb "github.com/gw-gong/gwkit-go/examples/case002/internal/protobuf"
)

func NewTestService() pb.TestServiceServer {
	return &testService{}
}

type testService struct {
	pb.UnimplementedTestServiceServer
}

func (s *testService) TestFunc(ctx context.Context, request *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{ResponseMsg: "test"}, nil
}
