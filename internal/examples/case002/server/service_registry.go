package main

import "github.com/gw-gong/gwkit-go/grpc/consul_agent"

const (
	ServiceName = consul_agent.ServiceName("test_service")
	ServiceTag  = "test"
	ServiceDC   = ""

	ServerPort = 8081
)

func NewTestServiceRegistry(sn consul_agent.ServiceName) (consul_agent.Registry, error) {
	registry, err := consul_agent.NewRegistry(sn)
	if err != nil {
		return nil, err
	}
	return registry, nil
}
