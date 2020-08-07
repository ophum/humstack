package vxlan

import (
	"log"
	"net"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestVxlanAdd(t *testing.T) {
	dev, err := netlink.LinkByName("enp0s31f6")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(dev.Attrs().Index)

	group := net.ParseIP("239.0.1.1")
	log.Println([]byte(group.To4()))
	_, err = NewVxlan("test-vxlan", 10, group, dev.Attrs().Index)
	if err != nil {
		t.Fatal(err)
	}

}
