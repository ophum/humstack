package network

import (
	"fmt"
	"hash/crc32"
	"log"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type NetworkAgent struct {
	client *client.Clients
}

const (
	NetworkV0AnnotationNetworkType = "networkv0/network-type"
)

const (
	NetworkV0NetworkTypeVXLAN  = "VXLAN"
	NetworkV0NetworkTypeVLAN   = "VLAN"
	NetworkV0NetworkTypeBridge = "Bridge"
)

func NewNetworkAgent(client *client.Clients) *NetworkAgent {
	return &NetworkAgent{
		client: client,
	}
}

func (a *NetworkAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			nsList, err := a.client.CoreV0().Namespace().List()
			if err != nil {
				continue
			}

			for _, ns := range nsList {
				netList, err := a.client.SystemV0().Network().List(ns.ID)
				if err != nil {
					continue
				}

				for _, net := range netList {
					switch net.Annotations[NetworkV0AnnotationNetworkType] {
					case NetworkV0NetworkTypeBridge:

						err = syncBridgeNetwork(net)
						if err != nil {
							continue
						}

						_, err = a.client.SystemV0().Network().Update(net)
						if err != nil {
							continue
						}
					}
				}
			}
		}
	}
}

func syncBridgeNetwork(network *system.Network) error {
	cs := crc32.Checksum([]byte(network.ID), crc32.IEEETable)

	bridgeName := fmt.Sprintf("hum-%010x", cs)
	log.Printf("create bridge `%s`\n", bridgeName)
	_, err := iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}

	network.Annotations["bridge_name"] = bridgeName
	return nil
}

func createVXLANNetwork(network *system.Network) error {

	return nil
}
