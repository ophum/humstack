package node

import (
	"log"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type NodeAgent struct {
	client   *client.Clients
	NodeInfo *system.Node
}

func NewNodeAgent(node system.Node, client *client.Clients) *NodeAgent {
	return &NodeAgent{
		NodeInfo: &node,
		client:   client,
	}
}

func (a *NodeAgent) Run() {

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			node, err := a.client.SystemV0().Node().Get(a.NodeInfo.Name)
			if err != nil {
				log.Println(err)
				continue
			}

			if node.Name == "" {
				node, err = a.client.SystemV0().Node().Create(a.NodeInfo)
				if err != nil {
					log.Println(err)
					continue
				}

				a.NodeInfo = node
			}

			if node.Status.State == system.NodeStateNotReady {
				node.Status.State = system.NodeStateReady
				node, err = a.client.SystemV0().Node().Update(node)
				if err != nil {
					log.Println(err)
					continue
				}
			}

			a.NodeInfo = node
		}
	}

}
