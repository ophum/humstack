package namespace

import (
	"time"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
)

type NamespaceAgent struct {
	client *client.Clients
	logger *zap.Logger
}

func NewNamespaceAgent(client *client.Clients, logger *zap.Logger) *NamespaceAgent {
	return &NamespaceAgent{
		client: client,
		logger: logger,
	}
}

func (a *NamespaceAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				a.logger.Error(
					"get group list",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get namespace list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
				}

				for _, ns := range nsList {
					if ns.DeleteState != meta.DeleteStateDelete {
						continue
					}
					// namespaceに所属するリソースにDeleteStateをセットする
					// virtualmachine, blockstorage, network, virtualrouter

					isDeletable := true
					if n, err := a.virtualMachinesSetDeleteState(ns); err != nil {
						a.logger.Error(
							"set virtualmachines delete state",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					} else if n != 0 {
						isDeletable = false
					}
					if n, err := a.blockStoragesSetDeleteState(ns); err != nil {
						a.logger.Error(
							"set blockstorages delete state",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					} else if n != 0 {
						isDeletable = false
					}
					if n, err := a.networksSetDeleteState(ns); err != nil {
						a.logger.Error(
							"set network delete state",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					} else if n != 0 {
						isDeletable = false
					}
					if n, err := a.nodeNetworksSetDeleteState(ns); err != nil {
						a.logger.Error(
							"set node network delete state",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					} else if n != 0 {
						isDeletable = false
					}
					if n, err := a.virtualRoutersSetDeleteState(ns); err != nil {
						a.logger.Error(
							"set virtualrouter delete state",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
					} else if n != 0 {
						isDeletable = false
					}

					if isDeletable {
						if err := a.client.CoreV0().Namespace().Delete(ns.Group, ns.ID); err != nil {
							a.logger.Error(
								"delete namespace",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
						}
					}
				}
			}
		}
	}
}

func (a *NamespaceAgent) virtualMachinesSetDeleteState(ns *core.Namespace) (int, error) {
	vmList, err := a.client.SystemV0().VirtualMachine().List(ns.Group, ns.ID)
	if err != nil {
		return -1, err
	}

	for _, vm := range vmList {
		if vm.DeleteState == meta.DeleteStateDelete {
			continue
		}

		_ = a.client.SystemV0().VirtualMachine().DeleteState(vm.Group, vm.Namespace, vm.ID)
	}

	return len(vmList), nil
}

func (a *NamespaceAgent) blockStoragesSetDeleteState(ns *core.Namespace) (int, error) {
	bsList, err := a.client.SystemV0().BlockStorage().List(ns.Group, ns.ID)
	if err != nil {
		return -1, err
	}

	for _, bs := range bsList {
		if bs.DeleteState == meta.DeleteStateDelete {
			continue
		}

		_ = a.client.SystemV0().BlockStorage().DeleteState(bs.Group, bs.Namespace, bs.ID)
	}

	return len(bsList), nil
}

func (a *NamespaceAgent) networksSetDeleteState(ns *core.Namespace) (int, error) {
	netList, err := a.client.CoreV0().Network().List(ns.Group, ns.ID)
	if err != nil {
		return -1, err
	}

	for _, net := range netList {
		if net.DeleteState == meta.DeleteStateDelete {
			continue
		}

		_ = a.client.CoreV0().Network().DeleteState(net.Group, net.Namespace, net.ID)
	}

	return len(netList), nil
}

func (a *NamespaceAgent) nodeNetworksSetDeleteState(ns *core.Namespace) (int, error) {
	netList, err := a.client.SystemV0().NodeNetwork().List(ns.Group, ns.ID)
	if err != nil {
		return -1, err
	}

	for _, net := range netList {
		if net.DeleteState == meta.DeleteStateDelete {
			continue
		}

		_ = a.client.SystemV0().NodeNetwork().DeleteState(net.Group, net.Namespace, net.ID)
	}

	return len(netList), nil
}
func (a *NamespaceAgent) virtualRoutersSetDeleteState(ns *core.Namespace) (int, error) {
	vrList, err := a.client.SystemV0().VirtualRouter().List(ns.Group, ns.ID)
	if err != nil {
		return -1, err
	}

	for _, vr := range vrList {
		if vr.DeleteState == meta.DeleteStateDelete {
			continue
		}

		_ = a.client.SystemV0().VirtualRouter().DeleteState(vr.Group, vr.Namespace, vr.ID)
	}

	return len(vrList), nil
}
