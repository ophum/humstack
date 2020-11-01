package node

import (
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
)

type NodeAgent struct {
	client   *client.Clients
	NodeInfo *system.Node
	logger   *zap.Logger
}

func NewNodeAgent(node system.Node, client *client.Clients, logger *zap.Logger) *NodeAgent {
	return &NodeAgent{
		NodeInfo: &node,
		client:   client,
		logger:   logger,
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
				a.logger.Error(
					"get node",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
				continue
			}

			if node.Name == "" {
				node, err = a.client.SystemV0().Node().Create(a.NodeInfo)
				if err != nil {
					a.logger.Error(
						"create node",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}

				a.NodeInfo = node
			}

			if node.Status.State == system.NodeStateNotReady {
				node.Status.State = system.NodeStateReady
				node, err = a.client.SystemV0().Node().Update(node)
				if err != nil {
					a.logger.Error(
						"update node",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}
			}

			a.NodeInfo = node
		}
	}

}

func (a *NodeAgent) GetNodeInfo() *system.Node {
	return a.NodeInfo
}
