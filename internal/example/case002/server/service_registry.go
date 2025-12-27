package main

import "github.com/gw-gong/gwkit-go/grpc/consul"

const (
	ServiceName = consul.ServiceName("test_service")
	ServiceTag  = "test"

	ServerPort = 8081
)

func NewTestServiceRegistry(sn consul.ServiceName) (consul.ConsulRegistry, error) {
	registry, err := consul.NewConsulRegistry(sn)
	if err != nil {
		return nil, err
	}
	return registry, nil
}
