package client

import "github.com/ophum/humstack/pkg/client/system"

type Clients struct {
	systemV0         *system.SystemV0Clients
	apiServerAddress string
	apiServerPort    int32
}

func NewClients(apiServerAddress string, apiServerPort int32) *Clients {
	return &Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		systemV0: system.NewSystemV0Clients(apiServerAddress, apiServerPort),
	}
}

func (c *Clients) SystemV0() *system.SystemV0Clients {
	return c.systemV0
}
