package system

import (
	nodev0 "github.com/ophum/humstack/pkg/client/system/node/v0"
)

type SystemV0Clients struct {
	apiServerAddress string
	apiServerPort    int32

	nodeClient *nodev0.NodeClient
}

func NewSystemV0Clients(apiServerAddress string, apiServerPort int32) *SystemV0Clients {
	nodeClient := nodev0.NewNodeClient("http", apiServerAddress, apiServerPort)
	return &SystemV0Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		nodeClient: nodeClient,
	}
}

func (c *SystemV0Clients) Node() *nodev0.NodeClient {
	return c.nodeClient
}
