package node

import (
	"fmt"
	"strconv"
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

func NewNodeAgent(node *system.Node, client *client.Clients, logger *zap.Logger) *NodeAgent {
	return &NodeAgent{
		NodeInfo: node,
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

			if node.ID == "" {
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

			if node.Status.State == system.NodeStateNotReady ||
				node.Status.State == "" {
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

			res, err := a.getUsedResources()
			if err != nil {
				a.logger.Error(
					"get requested resource",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
			} else {
				if node.Status.RequestedVcpus != res[ResourceTypeRequestVcpus] ||
					node.Status.RequestedMemory != res[ResourceTypeRequestMemory] {
					node.Status.RequestedVcpus = res[ResourceTypeRequestVcpus]
					node.Status.RequestedMemory = res[ResourceTypeRequestMemory]

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
			}

			a.NodeInfo = node
		}
	}

}

func (a *NodeAgent) GetNodeInfo() *system.Node {
	return a.NodeInfo
}

type ResourceType string

const (
	ResourceTypeRequestVcpus  ResourceType = "requestVcpus"
	ResourceTypeRequestMemory ResourceType = "requestMemory"
	ResourceTypeLimitVcpus    ResourceType = "limitVcpus"
	ResourceTypeLimitMemory   ResourceType = "limitMemory"
)

func (a *NodeAgent) getUsedResources() (map[ResourceType]string, error) {

	grList, err := a.client.CoreV0().Group().List()
	if err != nil {
		return nil, err
	}

	var vcpusRequests int64 = 0
	var vcpusLimits int64 = 0
	var memoryRequests int64 = 0
	var memoryLimits int64 = 0

	for _, group := range grList {
		nsList, err := a.client.CoreV0().Namespace().List(group.ID)
		if err != nil {
			return nil, err
		}

		for _, ns := range nsList {
			vmList, err := a.client.SystemV0().VirtualMachine().List(group.ID, ns.ID)
			if err != nil {
				return nil, err
			}

			for _, vm := range vmList {
				if vm.Annotations["virtualmachinev0/node_name"] != a.NodeInfo.ID {
					continue
				}
				if vm.Spec.ActionState == system.VirtualMachineActionStatePowerOff {
					continue
				}

				vcpusRequest, err := strconv.ParseInt(withUnitToWithoutUnit(vm.Spec.RequestVcpus), 10, 64)
				if err != nil {
					return nil, err
				}
				vcpusRequests += vcpusRequest

				vcpusLimit, err := strconv.ParseInt(withUnitToWithoutUnit(vm.Spec.LimitVcpus), 10, 64)
				if err != nil {
					return nil, err
				}
				vcpusLimits += vcpusLimit

				memoryRequest, err := strconv.ParseInt(withUnitToWithoutUnit(vm.Spec.RequestMemory), 10, 64)
				if err != nil {
					return nil, err
				}
				memoryRequests += memoryRequest

				memoryLimit, err := strconv.ParseInt(withUnitToWithoutUnit(vm.Spec.LimitMemory), 10, 64)
				if err != nil {
					return nil, err
				}
				memoryLimits += memoryLimit
			}
		}
	}

	return map[ResourceType]string{
		ResourceTypeRequestVcpus:  parseVcpus(vcpusRequests),
		ResourceTypeRequestMemory: parseMemory(memoryRequests),
		ResourceTypeLimitVcpus:    parseVcpus(vcpusLimits),
		ResourceTypeLimitMemory:   parseMemory(memoryLimits),
	}, nil
}

func parseVcpus(n int64) string {
	return fmt.Sprintf("%dm", n*1000)
}

func parseMemory(b int64) string {
	if b%(1024*1024*1024) == 0 {
		return fmt.Sprintf("%dG", b/(1024*1024*1024))
	} else if b%(1024*1024) == 0 {
		return fmt.Sprintf("%dM", b/(1024*1024))
	} else if b%1024 == 0 {
		return fmt.Sprintf("%dK", b/1024)
	} else {
		return fmt.Sprintf("%d", b)
	}
}

// copy from virtualmachine/agent.go
const (
	UnitGigabyte = 'G'
	UnitMegabyte = 'M'
	UnitKilobyte = 'K'
	UnitMilli    = 'm'
)

func withUnitToWithoutUnit(numberWithUnit string) string {
	length := len(numberWithUnit)
	if numberWithUnit[length-1] >= '0' && numberWithUnit[length-1] <= '9' {
		return numberWithUnit
	}

	number, err := strconv.ParseInt(numberWithUnit[:length-1], 10, 64)
	if err != nil {
		return "0"
	}

	switch numberWithUnit[length-1] {
	case UnitGigabyte:
		return fmt.Sprintf("%d", number*1024*1024*1024)
	case UnitMegabyte:
		return fmt.Sprintf("%d", number*1024*1024)
	case UnitKilobyte:
		return fmt.Sprintf("%d", number*1024)
	case UnitMilli:
		return fmt.Sprintf("%d", number/1000)
	}
	return "0"
}
