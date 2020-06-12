package veth

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

type Veth struct {
	Name     string
	PeerName string
	link     netlink.Link
	peerLink netlink.Link
}

func NewVeth(name, peerName string) (*Veth, error) {
	veth := &Veth{
		Name:     name,
		PeerName: peerName,
	}

	var err error
	veth.link, err = netlink.LinkByName(name)
	if err != nil {
		if err = veth.createVeth(); err != nil {
			return nil, err
		}
	} else {
		veth.peerLink, err = netlink.LinkByName(peerName)
		if err != nil {
			return nil, err
		}
	}

	if veth.link.Type() != "veth" {
		return nil, fmt.Errorf("Error: Interface is not veth")
	}

	return veth, nil
}

func (v *Veth) createVeth() error {
	la := netlink.NewLinkAttrs()
	la.Name = v.Name
	v.link = &netlink.Veth{
		LinkAttrs: la,
		PeerName:  v.PeerName,
	}

	if err := netlink.LinkAdd(v.link); err != nil {
		return fmt.Errorf("Error: failed to add link")
	}

	var err error
	v.peerLink, err = netlink.LinkByName(v.PeerName)
	if err != nil {
		return err
	}

	return v.Up()
}

func (v *Veth) Up() error {
	err := netlink.LinkSetUp(v.peerLink)
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(v.link)
}

func (v *Veth) Down() error {
	return netlink.LinkSetDown(v.link)
}

func (v *Veth) SetMaster(bridge *netlink.Bridge) error {
	return netlink.LinkSetMaster(v.link, bridge)
}

func (v *Veth) SetMasterPeer(bridge *netlink.Bridge) error {
	return netlink.LinkSetMaster(v.peerLink, bridge)
}

func (v *Veth) Delete() error {
	if err := netlink.LinkDel(v.link); err != nil {
		return err
	}

	v.link = nil
	return nil
}
