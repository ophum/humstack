package veth

import (
	"testing"

	"github.com/vishvananda/netlink"
)

func TestVethAdd(t *testing.T) {
	v, err := NewVeth("test", "test_peer")
	if err != nil {
		t.Fatal(err)
	}

	l, err := netlink.LinkByName("hum-br-055d73a0")
	if err != nil {
		t.Fatal(err)
	}

	err = v.SetMaster(l.(*netlink.Bridge))
	if err != nil {
		t.Fatal(err)
	}

}
