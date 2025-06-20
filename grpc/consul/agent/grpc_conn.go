package agent

import (
	"fmt"

	consul_api "github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

type ConsulServiceStatus int

const (
	ConsulServiceStatusUnknown ConsulServiceStatus = iota
	ConsulServiceStatusNotExists
	ConsulServiceStatusExists
)

func NewHealthyGrpcConn(agentAddr, serviceName, tag string, dc string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	serviceStatus := CheckServiceExists(agentAddr, serviceName, tag, true)
	switch serviceStatus {
	case ConsulServiceStatusUnknown:
		err = fmt.Errorf("failed to check service %s with tag %s", serviceName, tag)
		return nil, err
	case ConsulServiceStatusExists:
		return newGrpcConn(agentAddr, serviceName, tag, dc, opts...)
	case ConsulServiceStatusNotExists:
		err = fmt.Errorf("healthy service %s with tag %s not found", serviceName, tag)
		return nil, err
	}
	return nil, fmt.Errorf("unknown service status: %d", serviceStatus)
}

func CheckServiceExists(agentAddr, serviceName, tag string, passingOnly bool) ConsulServiceStatus {
	config := consul_api.DefaultConfig()
	config.Address = agentAddr
	client, err := consul_api.NewClient(config)
	if err != nil {
		return ConsulServiceStatusUnknown
	}

	entries, _, err := client.Health().Service(serviceName, tag, passingOnly, nil)
	if err != nil {
		return ConsulServiceStatusUnknown
	}
	if len(entries) > 0 {
		return ConsulServiceStatusExists
	}
	return ConsulServiceStatusNotExists
}

func newGrpcConn(agentAddr, serviceName, tag string, dc string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	target := formatGrpcConnTarget(agentAddr, serviceName, tag, dc)
	conn, err = grpc.NewClient(target, opts...)
	return conn, err
}

func formatGrpcConnTarget(agentAddr, serviceName, tag string, dc string) string {
	target := fmt.Sprintf("consul://%s/%s?healthy=true&tag=%s", agentAddr, serviceName, tag)
	if dc != "" {
		target += fmt.Sprintf("&dc=%s", dc)
	}
	return target
}
