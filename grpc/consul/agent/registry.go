package agent

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

// DefaultConsulAgentAddr is the default address of consul agent,
// seacloud applications may use a different value
const DefaultConsulAgentAddr = "127.0.0.1:8500"

type ServiceName string

// Registry defines interface for service registration and discovery
type Registry interface {
	// Get healthy service instance endpoints
	Service(tags []string) ([]string, error)
	// Register a service instance
	Register(serviceID string, port int, tags []string) error
	// Deregister a service instance
	Deregister(serviceID string) error
}

type registry struct {
	client *consul.Client
	name   string
}

// NewRegistry returns a Registry interface for services
func NewRegistry(name ServiceName) (*registry, error) {
	config := consul.DefaultConfig()
	c, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &registry{client: c, name: string(name)}, nil
}

func (r *registry) Service(tags []string) ([]string, error) {
	passingOnly := true
	addrs, _, err := r.client.Health().ServiceMultipleTags(r.name, tags, passingOnly, nil)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, fmt.Errorf("service ( %s ) was not found", r.name)
	}
	var endpoints = make([]string, 0, len(addrs))
	for _, addr := range addrs {
		endpoints = append(endpoints, fmt.Sprintf("%s:%d", addr.Node.Address, addr.Service.Port))
	}
	return endpoints, nil
}

func (r *registry) Register(serviceID string, port int, tags []string) error {
	reg := &consul.AgentServiceRegistration{
		ID:   serviceID,
		Name: r.name,
		Port: port,
		Tags: tags,
		Check: &consul.AgentServiceCheck{
			GRPC:       fmt.Sprintf("127.0.0.1:%d", port),
			GRPCUseTLS: false,
			Interval:   "10s",
			Timeout:    "3s",
			Name:       r.name + " grpc health check",
			Status:     consul.HealthPassing,
		},
	}
	return r.client.Agent().ServiceRegister(reg)
}

func (r *registry) RegisterTLS(serviceID string, port int, tags []string) error {
	reg := &consul.AgentServiceRegistration{
		ID:   serviceID,
		Name: r.name,
		Port: port,
		Tags: tags,
		Check: &consul.AgentServiceCheck{
			GRPC:          fmt.Sprintf("127.0.0.1:%d", port),
			GRPCUseTLS:    true,
			TLSSkipVerify: true,
			Interval:      "10s",
			Timeout:       "3s",
			Name:          r.name + " grpc health check",
			Status:        consul.HealthPassing,
		},
	}
	return r.client.Agent().ServiceRegister(reg)
}

func (r *registry) Deregister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}
