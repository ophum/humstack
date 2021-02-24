package nodenetwork

import (
	"fmt"
	"strconv"
	"time"

	"github.com/n0stack/n0stack/n0core/pkg/driver/iproute2"
	"github.com/ophum/humstack/pkg/agents/system/nodenetwork/utils"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

func (a *NodeNetworkAgent) syncVLANNetwork(network *system.NodeNetwork) error {

	bridgeName := utils.GenerateName("hum-br-", network.Group+network.Namespace+network.ID)
	vlanName := a.config.VLAN.DevName + "." + network.Spec.ID
	if a.config.VLAN.VLANInterfaceNamePrefix != "" {
		vlanName = a.config.VLAN.VLANInterfaceNamePrefix + "." + network.Spec.ID
	}

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
				network.Status.Logs = append(network.Status.Logs, system.NodeNetworkStatusLog{
					NodeID:   a.node,
					Datetime: time.Now().String(),
					Log:      fmt.Sprintf("vlan id `%s` is already used.", network.Spec.ID),
				})
				if _, err := a.client.SystemV0().NodeNetwork().Update(network); err != nil {
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

	// 削除処理
	if network.DeleteState == meta.DeleteStateDelete {
		if err := vlan.Delete(); err != nil {
			return errors.Wrap(err, "delete vlan")
		}

		if err := br.Delete(); err != nil {
			return errors.Wrap(err, "delete bridge")
		}
		if err := a.client.SystemV0().NodeNetwork().Delete(network.Group, network.Namespace, network.ID); err != nil {
			return errors.Wrap(err, "delete node network")
		}
		return nil
	}

	err = vlan.SetMaster(br)
	if err != nil {
		return err
	}

	network.Annotations[NetworkV0AnnotationBridgeName] = bridgeName
	network.Annotations[NetworkV0AnnotationVLANName] = vlanName
	network.Status.State = system.NetworkStateAvailable

	return setHash(network)
}
