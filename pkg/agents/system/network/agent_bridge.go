package network

import (
	"log"
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iptables"
	"github.com/ophum/humstack/pkg/agents/system/network/utils"
	"github.com/ophum/humstack/pkg/api/system"
)

const (
	NetworkV0AnnotationBridgeName     = "networkv0/bridge_name"
	NetworkV0AnnotationDefaultGateway = "networkv0/default_gateway"
)

func syncBridgeNetwork(network *system.Network) error {
	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	log.Printf("create bridge `%s`\n", bridgeName)
	br, err := iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}

	if gw, ok := network.Annotations[NetworkV0AnnotationDefaultGateway]; ok {
		_, ipnet, err := net.ParseCIDR(gw)
		if err != nil {
			return err
		}

		if err := br.SetAddress(gw); err != nil {
			return err
		}
		if err := iptables.CreateMasqueradeRule(bridgeName, ipnet); err != nil {
			return err
		}
	}

	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	return setHash(network)

}
