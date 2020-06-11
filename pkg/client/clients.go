package client

import (
	"github.com/ophum/humstack/pkg/client/core"
	"github.com/ophum/humstack/pkg/client/system"
)

type Clients struct {
	coreV0           *core.CoreV0Clients
	systemV0         *system.SystemV0Clients
	apiServerAddress string
	apiServerPort    int32
}

func NewClients(apiServerAddress string, apiServerPort int32) *Clients {
	return &Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		coreV0:   core.NewCoreV0Clients(apiServerAddress, apiServerPort),
		systemV0: system.NewSystemV0Clients(apiServerAddress, apiServerPort),
	}
}

func (c *Clients) CoreV0() *core.CoreV0Clients {
	return c.coreV0
}
func (c *Clients) SystemV0() *system.SystemV0Clients {
	return c.systemV0
}
