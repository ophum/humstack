package main

import (
	"github.com/ophum/humstack/pkg/agents/system/network"
	"github.com/ophum/humstack/pkg/agents/system/node"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

func main() {
	client := client.NewClients("localhost", 8080)
	nodeAgent := node.NewNodeAgent(system.Node{
		Meta: meta.Meta{
			Name: "test",
		},
		Spec: system.NodeSpec{
			LimitMemory: "8Gi",
			LimitVcpus:  "10000m",
		},
	}, client)

	netAgent := network.NewNetworkAgent(client)

	go nodeAgent.Run()
	netAgent.Run()

}
