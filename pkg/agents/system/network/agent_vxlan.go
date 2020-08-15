package network

import (
	"log"
	"net"
	"strconv"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/ophum/humstack/pkg/agents/system/network/utils"
	"github.com/ophum/humstack/pkg/agents/system/network/vxlan"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/vishvananda/netlink"
)

func (a *NetworkAgent) syncVXLANNetwork(network *system.Network) error {

	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	vxlanName := utils.GenerateName("hum-vx-", network.Group+network.Namespace+network.ID)
	log.Printf("create vxlan %s and bridge %s\n", vxlanName, bridgeName)

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
	err = vx.SetMaster(brLink.(*netlink.Bridge))
	if err != nil {
		return err
	}

	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	network.Annotations[NetworkV0AnnotationVXLANName] = vxlanName

	return setHash(network)
}
