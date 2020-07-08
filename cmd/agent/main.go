package main

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/agents/system/blockstorage"
	"github.com/ophum/humstack/pkg/agents/system/network"
	"github.com/ophum/humstack/pkg/agents/system/node"
	"github.com/ophum/humstack/pkg/agents/system/virtualmachine"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

func main() {
	client := client.NewClients("localhost", 8080)
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	nodeAgent := node.NewNodeAgent(system.Node{
		Meta: meta.Meta{
			Name: hostname,
		},
		Spec: system.NodeSpec{
			LimitMemory: "8Gi",
			LimitVcpus:  "10000m",
		},
	}, client)

	netAgent := network.NewNetworkAgent(client)

	bsAgent := blockstorage.NewBlockStorageAgent(client, "./blockstorages")

	vmAgent := virtualmachine.NewVirtualMachineAgent(client)

	go nodeAgent.Run()
	go bsAgent.Run()
	go vmAgent.Run()
	netAgent.Run()

}
