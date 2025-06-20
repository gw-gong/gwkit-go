package main

import "github.com/gw-gong/gwkit-go/grpc/consul/agent"

const (
	ServiceName = agent.ServiceName("test_service")
	ServiceTag  = "test"
	ServiceDC   = ""

	ServerPort = 8080
)

func NewTestServiceRegistry(sn agent.ServiceName) (agent.Registry, error) {
	registry, err := agent.NewRegistry(sn)
	if err != nil {
		return nil, err
	}
	return registry, nil
}
