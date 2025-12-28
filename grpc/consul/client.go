package consul

import (
	"fmt"

	"github.com/gw-gong/gwkit-go/util/str"

	consul_api "github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

type AgentAddr string

const DefaultConsulAgentAddr AgentAddr = "127.0.0.1:8500"

type RegisterEntry struct {
	ServiceName string   `json:"service_name" yaml:"service_name" mapstructure:"service_name"`
	ServiceID   string   `json:"service_id" yaml:"service_id" mapstructure:"service_id"`
	Tags        []string `json:"tags" yaml:"tags" mapstructure:"tags"`
}

func (r *RegisterEntry) GenerateServiceID() string {
	return str.GenerateUUID()
}

type HealthyGrpcConnEntry struct {
	ServiceName string            `json:"service_name" yaml:"service_name" mapstructure:"service_name"`
	Tag         string            `json:"tag" yaml:"tag" mapstructure:"tag"`
	Opts        []grpc.DialOption `json:"-" yaml:"-" mapstructure:"-"`
}

type ConsulServiceStatus int

const (
	ConsulServiceStatusUnknown ConsulServiceStatus = iota
	ConsulServiceStatusNotExists
	ConsulServiceStatusExists
)

type ConsulClient interface {
	// Register a service instance
	Register(entry *RegisterEntry, port int, useTLS bool) error
	// Deregister a service instance
	Deregister(serviceID string) error
	// GetHealthyGrpcConn returns a healthy gRPC connection
	GetHealthyGrpcConn(entry *HealthyGrpcConnEntry) (conn *grpc.ClientConn, err error)
}

type consulClient struct {
	client    *consul_api.Client
	agentAddr string
}

// NewConsulRegistry returns a Registry interface for services
func NewConsulClient(agentAddr AgentAddr) (ConsulClient, error) {
	if agentAddr == "" {
		return nil, fmt.Errorf("agent address is required")
	}
	config := consul_api.DefaultConfig()
	config.Address = string(agentAddr)
	c, err := consul_api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}
	return &consulClient{client: c, agentAddr: string(agentAddr)}, nil
}

func (r *consulClient) Register(entry *RegisterEntry, port int, useTLS bool) error {
	check := &consul_api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", "127.0.0.1", port),
		Interval: "10s",
		Timeout:  "3s",
		Name:     entry.ServiceName + " grpc health check",
		Status:   consul_api.HealthPassing,
	}

	if useTLS {
		check.GRPCUseTLS = true
		check.TLSSkipVerify = true
	}

	reg := &consul_api.AgentServiceRegistration{
		ID:    entry.ServiceID,
		Name:  entry.ServiceName,
		Port:  port,
		Tags:  entry.Tags,
		Check: check,
	}

	err := r.client.Agent().ServiceRegister(reg)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	return nil
}

func (r *consulClient) Deregister(serviceID string) error {
	err := r.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	return nil
}

func (r *consulClient) GetHealthyGrpcConn(entry *HealthyGrpcConnEntry) (conn *grpc.ClientConn, err error) {
	if entry == nil {
		return nil, fmt.Errorf("entry is nil")
	}
	if entry.ServiceName == "" || entry.Tag == "" {
		return nil, fmt.Errorf("service name, and tag are required")
	}
	serviceStatus := r.checkServiceExists(entry.ServiceName, entry.Tag, true)
	switch serviceStatus {
	case ConsulServiceStatusUnknown:
		return nil, fmt.Errorf("failed to check service %s with tag %s: %w", entry.ServiceName, entry.Tag, err)
	case ConsulServiceStatusExists:
		return r.newGrpcConn(entry.ServiceName, entry.Tag, entry.Opts...)
	case ConsulServiceStatusNotExists:
		return nil, fmt.Errorf("healthy service %s with tag %s not found: %w", entry.ServiceName, entry.Tag, err)
	}
	return nil, fmt.Errorf("unknown service status: %d", serviceStatus)
}

func (r *consulClient) checkServiceExists(serviceName, tag string, passingOnly bool) ConsulServiceStatus {
	entries, _, err := r.client.Health().Service(serviceName, tag, passingOnly, nil)
	if err != nil {
		return ConsulServiceStatusUnknown
	}
	if len(entries) > 0 {
		return ConsulServiceStatusExists
	}
	return ConsulServiceStatusNotExists
}

func (r *consulClient) newGrpcConn(serviceName, tag string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	target := formatGrpcConnTarget(r.agentAddr, serviceName, tag)
	conn, err = grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc connection: %w", err)
	}
	return conn, nil
}

func formatGrpcConnTarget(agentAddr, serviceName, tag string) string {
	return fmt.Sprintf("consul://%s/%s?healthy=true&tag=%s", agentAddr, serviceName, tag)
}
