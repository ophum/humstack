package network

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

type NetworkAgent struct {
	client *client.Clients
	config *NetworkAgentConfig
	node   string
	logger *zap.Logger
}

const (
	NetworkV0AnnotationNetworkType = "networkv0/network_type"
)

const (
	NetworkV0NetworkTypeVXLAN  = "VXLAN"
	NetworkV0NetworkTypeVLAN   = "VLAN"
	NetworkV0NetworkTypeBridge = "Bridge"
)

func NewNetworkAgent(client *client.Clients, config *NetworkAgentConfig, logger *zap.Logger) *NetworkAgent {
	node, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &NetworkAgent{
		client: client,
		config: config,
		node:   node,
		logger: logger,
	}
}

func (a *NetworkAgent) Run() {
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

					netList, err := a.client.SystemV0().Network().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get network list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					for _, net := range netList {
						oldHash := net.ResourceHash
						net.Status.AttachedInterfaces = attachedInterfacesToNet[net.ID]
						switch net.Annotations[NetworkV0AnnotationNetworkType] {
						case NetworkV0NetworkTypeBridge:
							err = syncBridgeNetwork(net)
							if err != nil {
								a.logger.Error(
									"sync bridge network",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}
						case NetworkV0NetworkTypeVXLAN:
							err = a.syncVXLANNetwork(net)
							if err != nil {
								a.logger.Error(
									"sync vxlan network",
									zap.String("msg", err.Error()),
									zap.Time("time", time.Now()),
								)
								continue
							}
						case NetworkV0NetworkTypeVLAN:
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
						_, err = a.client.SystemV0().Network().Update(net)
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

func setHash(network *system.Network) error {
	network.ResourceHash = ""
	resourceJSON, err := json.Marshal(network)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	network.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
