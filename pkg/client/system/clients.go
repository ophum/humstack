package system

import (
	netv0 "github.com/ophum/humstack/pkg/client/system/network/v0"
	nodev0 "github.com/ophum/humstack/pkg/client/system/node/v0"
)

type SystemV0Clients struct {
	apiServerAddress string
	apiServerPort    int32

	nodeClient    *nodev0.NodeClient
	networkClient *netv0.NetworkClient
}

func NewSystemV0Clients(apiServerAddress string, apiServerPort int32) *SystemV0Clients {
	nodeClient := nodev0.NewNodeClient("http", apiServerAddress, apiServerPort)
	return &SystemV0Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		nodeClient:    nodeClient,
		networkClient: netv0.NewNetworkClient("http", apiServerAddress, apiServerPort),
	}
}

func (c *SystemV0Clients) Node() *nodev0.NodeClient {
	return c.nodeClient
}

func (c *SystemV0Clients) Network() *netv0.NetworkClient {
	return c.networkClient
}
