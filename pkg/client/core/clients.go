package core

import (
	grv0 "github.com/ophum/humstack/pkg/client/core/group/v0"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
)

type CoreV0Clients struct {
	apiServerAddress string
	apiServerPort    int32

	namespaceClient *nsv0.NamespaceClient
	groupClient     *grv0.GroupClient
}

func NewCoreV0Clients(apiServerAddress string, apiServerPort int32) *CoreV0Clients {
	return &CoreV0Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		namespaceClient: nsv0.NewNamespaceClient("http", apiServerAddress, apiServerPort),
		groupClient:     grv0.NewGroupClient("http", apiServerAddress, apiServerPort),
	}
}

func (c *CoreV0Clients) Namespace() *nsv0.NamespaceClient {
	return c.namespaceClient
}

func (c *CoreV0Clients) Group() *grv0.GroupClient {
	return c.groupClient
}
