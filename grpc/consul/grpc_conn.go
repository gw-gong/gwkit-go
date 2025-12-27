package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

type ConsulServiceStatus int

const (
	ConsulServiceStatusUnknown ConsulServiceStatus = iota
	ConsulServiceStatusNotExists
	ConsulServiceStatusExists
)

type HealthyGrpcConnOption struct {
	AgentAddr   string            `json:"agent_addr" yaml:"agent_addr" mapstructure:"agent_addr"`
	ServiceName string            `json:"service_name" yaml:"service_name" mapstructure:"service_name"`
	Tag         string            `json:"tag" yaml:"tag" mapstructure:"tag"`
	Opts        []grpc.DialOption `json:"-" yaml:"-" mapstructure:"-"`
}

func NewHealthyGrpcConn(option *HealthyGrpcConnOption) (conn *grpc.ClientConn, err error) {
	if option == nil {
		return nil, fmt.Errorf("option is nil")
	}
	if option.AgentAddr == "" || option.ServiceName == "" || option.Tag == "" {
		return nil, fmt.Errorf("agent addr, service name, and tag are required")
	}
	serviceStatus := CheckServiceExists(option.AgentAddr, option.ServiceName, option.Tag, true)
	switch serviceStatus {
	case ConsulServiceStatusUnknown:
		err = fmt.Errorf("failed to check service %s with tag %s", option.ServiceName, option.Tag)
		return nil, err
	case ConsulServiceStatusExists:
		return newGrpcConn(option.AgentAddr, option.ServiceName, option.Tag, option.Opts...)
	case ConsulServiceStatusNotExists:
		err = fmt.Errorf("healthy service %s with tag %s not found", option.ServiceName, option.Tag)
		return nil, err
	}
	return nil, fmt.Errorf("unknown service status: %d", serviceStatus)
}

func CheckServiceExists(agentAddr, serviceName, tag string, passingOnly bool) ConsulServiceStatus {
	config := consulapi.DefaultConfig()
	config.Address = agentAddr
	client, err := consulapi.NewClient(config)
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

func newGrpcConn(agentAddr, serviceName, tag string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	target := formatGrpcConnTarget(agentAddr, serviceName, tag)
	conn, err = grpc.NewClient(target, opts...)
	return conn, err
}

func formatGrpcConnTarget(agentAddr, serviceName, tag string) string {
	return fmt.Sprintf("consul://%s/%s?healthy=true&tag=%s", agentAddr, serviceName, tag)
}
