package core

import (
	eipv0 "github.com/ophum/humstack/pkg/client/core/externalip/v0"
	eippoolv0 "github.com/ophum/humstack/pkg/client/core/externalippool/v0"
	grv0 "github.com/ophum/humstack/pkg/client/core/group/v0"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
	netv0 "github.com/ophum/humstack/pkg/client/core/network/v0"
)

type CoreV0Clients struct {
	apiServerAddress string
	apiServerPort    int32

	namespaceClient *nsv0.NamespaceClient
	groupClient     *grv0.GroupClient
	eippoolClient   *eippoolv0.ExternalIPPoolClient
	eipClient       *eipv0.ExternalIPClient
	networkClient   *netv0.NetworkClient
}

func NewCoreV0Clients(apiServerAddress string, apiServerPort int32) *CoreV0Clients {
	return &CoreV0Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		namespaceClient: nsv0.NewNamespaceClient("http", apiServerAddress, apiServerPort),
		groupClient:     grv0.NewGroupClient("http", apiServerAddress, apiServerPort),
		eipClient:       eipv0.NewExternalIPClient("http", apiServerAddress, apiServerPort),
		eippoolClient:   eippoolv0.NewExternalIPPoolClient("http", apiServerAddress, apiServerPort),
		networkClient:   netv0.NewNetworkClient("http", apiServerAddress, apiServerPort),
	}
}

func (c *CoreV0Clients) Namespace() *nsv0.NamespaceClient {
	return c.namespaceClient
}

func (c *CoreV0Clients) Group() *grv0.GroupClient {
	return c.groupClient
}

func (c *CoreV0Clients) ExternalIPPool() *eippoolv0.ExternalIPPoolClient {
	return c.eippoolClient
}

func (c *CoreV0Clients) ExternalIP() *eipv0.ExternalIPClient {
	return c.eipClient
}

func (c *CoreV0Clients) Network() *netv0.NetworkClient {
	return c.networkClient
}
