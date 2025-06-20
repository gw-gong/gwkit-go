package main

import "github.com/gw-gong/gwkit-go/grpc/consul/agent"

const (
	serviceName = agent.ServiceName("test_service")
	serviceTag  = "test"
	serviceDC   = ""

	serverPort = 8080
)

func NewTestServiceRegistry(sn agent.ServiceName) (agent.Registry, error) {
	registry, err := agent.NewRegistry(sn)
	if err != nil {
		return nil, err
	}
	return registry, nil
}
