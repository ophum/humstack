package nodenetwork

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"go.uber.org/zap"
)

type NodeNetworkAgent struct {
	client *client.Clients
	config *NetworkAgentConfig
	node   string
	logger *zap.Logger
}

const (
	NodeNetworkV0AnnotationNetworkType = "nodenetworkv0/network_type"
	NodeNetworkV0AnnotationNodeName    = "nodenetworkv0/node_name"
)

const (
	NodeNetworkV0NetworkTypeVXLAN  = "VXLAN"
	NodeNetworkV0NetworkTypeVLAN   = "VLAN"
	NodeNetworkV0NetworkTypeBridge = "Bridge"
)

func NewNodeNetworkAgent(client *client.Clients, config *NetworkAgentConfig, logger *zap.Logger) *NodeNetworkAgent {
	node, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &NodeNetworkAgent{
		client: client,
		config: config,
		node:   node,
		logger: logger,
	}
}

func (a *NodeNetworkAgent) Run(pollingDuration time.Duration) {
	ticker := time.NewTicker(pollingDuration)
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
				continue
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get namespace list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}

				for _, ns := range nsList {
					vmList, err := a.client.SystemV0().VirtualMachine().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get virtualmachine list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}
					attachedInterfacesToNet := map[string]map[string]system.VirtualMachineNIC{}
					for _, vm := range vmList {
						if vm.Status.State != system.VirtualMachineStateRunning {
							continue
						}

						for _, nic := range vm.Spec.NICs {
							if attachedInterfacesToNet[nic.NetworkID] == nil {
								attachedInterfacesToNet[nic.NetworkID] = map[string]system.VirtualMachineNIC{}
							}

							attachedInterfacesToNet[nic.NetworkID][filepath.Join("virtualmachinev0", vm.ID)] = *nic
						}
					}

					vrList, err := a.client.SystemV0().VirtualRouter().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get virtualrouter list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}
					for _, vr := range vrList {
						if vr.Status.State != system.VirtualRouterStateRunning {
							continue
						}

						for _, nic := range vr.Spec.NICs {
							if attachedInterfacesToNet[nic.NetworkID] == nil {
								attachedInterfacesToNet[nic.NetworkID] = map[string]system.VirtualMachineNIC{}
							}

							attachedInterfacesToNet[nic.NetworkID][filepath.Join("virtualrouterv0", vr.ID)] = system.VirtualMachineNIC{
								NetworkID:   nic.NetworkID,
								IPv4Address: nic.IPv4Address,
							}
						}
					}

					netList, err := a.client.SystemV0().NodeNetwork().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get network list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					for _, net := range netList {
						if nodeName, ok := net.Annotations[NodeNetworkV0AnnotationNodeName]; ok && nodeName != a.node {
							continue
						}
						oldHash := net.ResourceHash
						net.Status.AttachedInterfaces = attachedInterfacesToNet[net.ID]
						switch net.Annotations[NodeNetworkV0AnnotationNetworkType] {
						case NodeNetworkV0NetworkTypeBridge:
							err = a.syncBridgeNetwork(net)
							if err != nil {
								a.logger.Error(
									"sync bridge network",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}
						case NodeNetworkV0NetworkTypeVXLAN:
							err = a.syncVXLANNetwork(net)
							if err != nil {
								a.logger.Error(
									"sync vxlan network",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}
						case NodeNetworkV0NetworkTypeVLAN:
							err = a.syncVLANNetwork(net)
							if err != nil {
								a.logger.Error(
									"sync vlan network",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}

						}

						if net.ResourceHash == oldHash {
							continue
						}
						_, err = a.client.SystemV0().NodeNetwork().Update(net)
						if err != nil {
							a.logger.Error(
								"update network",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}
					}
				}
			}
		}
	}
}

func setHash(network *system.NodeNetwork) error {
	network.ResourceHash = ""
	resourceJSON, err := json.Marshal(network)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	network.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
