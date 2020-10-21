package network

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/ophum/humstack/pkg/agents/system/network/utils"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/vishvananda/netlink"
)

func (a *NetworkAgent) syncVLANNetwork(network *system.Network) error {

	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	vlanName := a.config.VLAN.DevName + "." + network.Spec.ID
	log.Printf("create vlan %s and bridge %s\n", vlanName, bridgeName)

	// vlan idがすでに別のbridgeに接続されているかチェックする
	vlanLink, err := netlink.LinkByName(vlanName)
	if err != nil {
		if err.Error() != "Link not found" {
			return err
		}
	}
	if err == nil {
		index := vlanLink.Attrs().MasterIndex
		attachedBr, err := netlink.LinkByIndex(index)
		if err != nil {
			if err.Error() != "Link not found" {
				return err
			}
		} else {
			if bridgeName != attachedBr.Attrs().Name {
				// vlan id is already used
				network.Status.Logs = append(network.Status.Logs, system.NetworkStatusLog{
					NodeID:   a.node,
					Datetime: time.Now().String(),
					Log:      fmt.Sprintf("vlan id `%s` is already used.", network.Spec.ID),
				})
				if _, err := a.client.SystemV0().Network().Update(network); err != nil {
					return err
				}
				return fmt.Errorf("vlan id `%s` is already used.", network.Spec.ID)
			}
		}
	}

	// 作成だけ
	_, err = iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(network.Spec.ID, 10, 64)
	if err != nil {
		return err
	}
	dev, err := iproute2.GetInterface(a.config.VLAN.DevName)
	if err != nil {
		return err
	}

	vlan, err := iproute2.NewVlan(dev, int(id))
	if err != nil {
		return err
	}

	br, err := iproute2.NewBridge(bridgeName)
	if err != nil {
		return err
	}
	err = vlan.SetMaster(br)
	if err != nil {
		return err
	}

	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	network.Annotations[NetworkV0AnnotationVLANName] = vlanName

	return setHash(network)
}
