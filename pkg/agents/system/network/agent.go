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
)

type NetworkAgent struct {
	client *client.Clients
	config *NetworkAgentConfig
	node   string
}

const (
	NetworkV0AnnotationNetworkType = "networkv0/network_type"
)

const (
	NetworkV0NetworkTypeVXLAN  = "VXLAN"
	NetworkV0NetworkTypeVLAN   = "VLAN"
	NetworkV0NetworkTypeBridge = "Bridge"
)

func NewNetworkAgent(client *client.Clients, config *NetworkAgentConfig) *NetworkAgent {
	node, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return &NetworkAgent{
		client: client,
		config: config,
		node:   node,
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
				log.Printf("[NET] %s", err.Error())
				continue
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					continue
				}

				for _, ns := range nsList {
					vmList, err := a.client.SystemV0().VirtualMachine().List(group.ID, ns.ID)
					if err != nil {
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
						continue
					}

					for _, net := range netList {
						oldHash := net.ResourceHash
						net.Status.AttachedInterfaces = attachedInterfacesToNet[net.ID]
						switch net.Annotations[NetworkV0AnnotationNetworkType] {
						case NetworkV0NetworkTypeBridge:
							err = syncBridgeNetwork(net)
							if err != nil {
								log.Println("error sync bridge network")
								log.Println(err)
								continue
							}
						case NetworkV0NetworkTypeVXLAN:
							err = a.syncVXLANNetwork(net)
							if err != nil {
								log.Println("error sync vxlan network")
								log.Println(err)
								continue
							}
						case NetworkV0NetworkTypeVLAN:
							err = a.syncVLANNetwork(net)
							if err != nil {
								log.Println("error sync vlan network")
								log.Println(err)
								continue
							}

						}

						if net.ResourceHash == oldHash {
							log.Println("no update")
							continue
						}
						log.Println("update store")
						_, err = a.client.SystemV0().Network().Update(net)
						if err != nil {
							log.Println("error Update store")
							log.Println(err)
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
