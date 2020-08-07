package vxlan

import (
	"net"

	"github.com/vishvananda/netlink"
)

type Vxlan struct {
	Name     string
	ID       int
	Group    net.IP
	DevIndex int
	link     netlink.Link
}

func NewVxlan(name string, id int, group net.IP, devIndex int) (*Vxlan, error) {
	vxlan := &Vxlan{
		Name:     name,
		ID:       id,
		Group:    group,
		DevIndex: devIndex,
	}

	var err error
	vxlan.link, err = netlink.LinkByName(name)
	if err != nil {
		if err = vxlan.createVxlan(); err != nil {
			return nil, err
		}
	}

	return vxlan, nil
}

func (v *Vxlan) createVxlan() error {
	la := netlink.NewLinkAttrs()
	la.Name = v.Name
	v.link = &netlink.Vxlan{
		LinkAttrs:    la,
		VxlanId:      v.ID,
		Group:        v.Group,
		VtepDevIndex: v.DevIndex,
		Learning:     true,
		L2miss:       true,
		L3miss:       true,
	}

	if err := netlink.LinkAdd(v.link); err != nil {
		return err
	}

	return v.Up()
}

func (v *Vxlan) Up() error {
	return netlink.LinkSetUp(v.link)
}

func (v *Vxlan) Down() error {
	return netlink.LinkSetDown(v.link)
}

func (v *Vxlan) SetMaster(bridge *netlink.Bridge) error {
	return netlink.LinkSetMaster(v.link, bridge)
}

func (v *Vxlan) Delete() error {
	if err := netlink.LinkDel(v.link); err != nil {
		return err
	}

	v.link = nil
	return nil
}
