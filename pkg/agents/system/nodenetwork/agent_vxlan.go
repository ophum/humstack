package nodenetwork

import (
	"net"
	"strconv"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/ophum/humstack/pkg/agents/system/nodenetwork/utils"
	"github.com/ophum/humstack/pkg/agents/system/nodenetwork/vxlan"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func (a *NodeNetworkAgent) syncVXLANNetwork(network *system.NodeNetwork) error {

	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	vxlanName := utils.GenerateName("hum-vx-", network.Group+network.Namespace+network.ID)

	// 作成だけ
	_, err := iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(network.Spec.ID, 10, 64)
	if err != nil {
		return err
	}
	dev, err := netlink.LinkByName(a.config.VXLAN.DevName)
	if err != nil {
		return err
	}
	vx, err := vxlan.NewVxlan(vxlanName, int(id), net.ParseIP(a.config.VXLAN.Group), dev.Attrs().Index)
	if err != nil {
		return err
	}

	brLink, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	if network.DeleteState == meta.DeleteStateDelete {
		br, err := iproute2.NewBridge(bridgeName)
		if err != nil {
			return err
		}
		if err := br.Delete(); err != nil {
			return errors.Wrap(err, "delete bridge")
		}

		if err := vx.Delete(); err != nil {
			return errors.Wrap(err, "delete vxlan")
		}

		if err := a.client.SystemV0().NodeNetwork().Delete(network.Group, network.Namespace, network.ID); err != nil {
			return errors.Wrap(err, "delete node network")
		}
		return nil
	}
	err = vx.SetMaster(brLink.(*netlink.Bridge))
	if err != nil {
		return err
	}

	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	network.Annotations[NetworkV0AnnotationVXLANName] = vxlanName

	return setHash(network)
}
