package server

import (
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	capi "github.com/hashicorp/consul/api"
	"micros-api/internal/conf"
)

type ConsulClient struct {
	Client *capi.Client
}

// NewConsulClient .
func NewConsulClient(c *conf.Data) (*ConsulClient, error) {
	config := capi.DefaultConfig()
	config.Address = c.Consul.Addr
	config.Token = c.Consul.Token
	client, err := capi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{
		Client: client,
	}, nil
}

// NewRegistry .
func NewRegistry(client *ConsulClient) *consul.Registry {
	return consul.New(client.Client)
}
