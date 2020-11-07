package nodenetwork

import (
	"net"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/n0stack/n0stack/n0core/pkg/driver/iptables"
	"github.com/ophum/humstack/pkg/agents/system/nodenetwork/utils"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
)

func (a *NodeNetworkAgent) syncBridgeNetwork(network *system.NodeNetwork) error {
	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	br, err := iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}

	if network.DeleteState == meta.DeleteStateDelete {
		if err := br.Delete(); err != nil {
			return errors.Wrap(err, "delete bridge")
		}

		if err := a.client.SystemV0().NodeNetwork().Delete(network.Group, network.Namespace, network.ID); err != nil {
			return errors.Wrap(err, "delete node network")
		}
		return nil
	}

	if gw, ok := network.Annotations[NetworkV0AnnotationDefaultGateway]; ok {
		if gw == "" {
			goto END
		}
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
END:
	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	return setHash(network)

}
