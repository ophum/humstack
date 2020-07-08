package system

import (
	bsv0 "github.com/ophum/humstack/pkg/client/system/blockstorage/v0"
	netv0 "github.com/ophum/humstack/pkg/client/system/network/v0"
	nodev0 "github.com/ophum/humstack/pkg/client/system/node/v0"
	vmv0 "github.com/ophum/humstack/pkg/client/system/virtualmachine/v0"
)

type SystemV0Clients struct {
	apiServerAddress string
	apiServerPort    int32

	nodeClient           *nodev0.NodeClient
	networkClient        *netv0.NetworkClient
	blockstorageClient   *bsv0.BlockStorageClient
	virtualmachineClient *vmv0.VirtualMachineClient
}

func NewSystemV0Clients(apiServerAddress string, apiServerPort int32) *SystemV0Clients {
	nodeClient := nodev0.NewNodeClient("http", apiServerAddress, apiServerPort)
	return &SystemV0Clients{
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,

		nodeClient:           nodeClient,
		networkClient:        netv0.NewNetworkClient("http", apiServerAddress, apiServerPort),
		blockstorageClient:   bsv0.NewBlockStorageClient("http", apiServerAddress, apiServerPort),
		virtualmachineClient: vmv0.NewVirtualMachineClient("http", apiServerAddress, apiServerPort),
	}
}

func (c *SystemV0Clients) Node() *nodev0.NodeClient {
	return c.nodeClient
}

func (c *SystemV0Clients) Network() *netv0.NetworkClient {
	return c.networkClient
}

func (c *SystemV0Clients) BlockStorage() *bsv0.BlockStorageClient {
	return c.blockstorageClient
}

func (c *SystemV0Clients) VirtualMachine() *vmv0.VirtualMachineClient {
	return c.virtualmachineClient
}
