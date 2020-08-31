package virtualrouter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ophum/humstack/pkg/agents/system/network/utils"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/vishvananda/netlink"
)

type VirtualRouterAgent struct {
	client *client.Clients

	externalBridge  string
	floatingIPCIDR  string
	usedFloatingIPs map[string]bool
}

const (
	VirtualRouterV0AnnotationNodeName = "virtualrouterv0/node_name"
)

func NewVirtualRouterAgent(client *client.Clients, externalBridge string, floatingIPCIDR string, usedFloatingIPs []string) *VirtualRouterAgent {
	return &VirtualRouterAgent{
		client:          client,
		externalBridge:  externalBridge,
		floatingIPCIDR:  floatingIPCIDR,
		usedFloatingIPs: map[string]bool{},
	}
}

func (a *VirtualRouterAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ticker.C:
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, ns := range nsList {
					vrList, err := a.client.SystemV0().VirtualRouter().List(group.ID, ns.ID)
					if err != nil {
						log.Println(err)
						continue
					}

					for _, vr := range vrList {
						oldHash := vr.ResourceHash
						if vr.Annotations[VirtualRouterV0AnnotationNodeName] != nodeName {
							log.Println("continue")
							continue
						}

						err = a.syncVirtualRouter(vr)
						if err != nil {
							log.Println(err)
							continue
						}

						if vr.ResourceHash == oldHash {
							log.Printf("vrouter(`%s`) no update\n", vr.ID)
							continue
						}

						_, err := a.client.SystemV0().VirtualRouter().Update(vr)
						if err != nil {
							log.Println(err)
							continue
						}
					}
				}
			}
		}
	}
}

func (a *VirtualRouterAgent) syncVirtualRouter(vr *system.VirtualRouter) error {

	log.Println("[VIRTUAL ROUTER]")
	netnsName := utils.GenerateName("netns-", vr.Group+vr.Namespace+vr.ID)
	rtExVeth := utils.GenerateName("hrt-ex-", netnsName)
	exRtVeth := utils.GenerateName("hex-rt-", netnsName)
	if !netnsIsExists(netnsName) {
		log.Println("[VR] netns add")
		if err := netnsAdd(netnsName); err != nil {
			return err
		}

		log.Println("[VR] link add")
		if err := ipLinkAddVeth(rtExVeth, exRtVeth); err != nil {
			return err
		}

		log.Println("[VR] set netns")
		if err := ipLinkSetNetNS(rtExVeth, netnsName); err != nil {
			return err
		}

		log.Println("[VR] set master")
		if err := ipLinkSetMaster(exRtVeth, a.externalBridge); err != nil {
			return err
		}

		log.Println("[VR] ip address add rtExVeth")
		netnsExec(netnsName, []string{
			"ip", "a", "add", vr.Spec.NATGatewayIP, "dev", rtExVeth,
		})
	}

	for _, eip := range vr.Spec.ExternalIPs {
		netnsExec(netnsName, []string{
			"ip", "a", "add", eip.IPv4Address, "dev", rtExVeth,
		})

		// iptables -t nat -A PREROUTING -d ${daddr} -j DNAT --to-destination ${DEST}
		// iptables -t nat -A POSTROUTING -s ${saddr} -j SNAT --to-source ${daddr}

		daddr := strings.Split(eip.IPv4Address, "/")[0]

		err := netnsExec(netnsName, []string{
			"iptables",
			"-t", "nat",
			"-C", "PREROUTING",
			"-d", daddr,
			"-j", "DNAT",
			"--to-destination", eip.BindInternalIPv4Address,
		})
		if err != nil {
			netnsExec(netnsName, []string{
				"iptables",
				"-t", "nat",
				"-A", "PREROUTING",
				"-d", daddr,
				"-j", "DNAT",
				"--to-destination", eip.BindInternalIPv4Address,
			})
		}
		err = netnsExec(netnsName, []string{
			"iptables",
			"-t", "nat",
			"-C", "POSTROUTING",
			"-s", eip.BindInternalIPv4Address,
			"-j", "SNAT",
			"--to-source", daddr,
		})
		if err != nil {
			log.Println("snat")
			err = netnsExec(netnsName, []string{
				"iptables",
				"-t", "nat",
				"-A", "POSTROUTING",
				"-s", eip.BindInternalIPv4Address,
				"-j", "SNAT",
				"--to-source", strings.Split(eip.IPv4Address, "/")[0],
			})
			log.Println(err)
		}
	}

	natGatewayIP := strings.Split(vr.Spec.NATGatewayIP, "/")[0]
	for _, nic := range vr.Spec.NICs {
		log.Println(nic.NetworkID)
		rtBrVeth := utils.GenerateName("hrt-br-", netnsName+nic.NetworkID)
		brRtVeth := utils.GenerateName("hbr-rt-", netnsName+nic.NetworkID)
		log.Println(vr.Group + vr.Namespace + nic.NetworkID)
		brName := utils.GenerateName("hum-br-", vr.Group+vr.Namespace+nic.NetworkID)
		_, err := netlink.LinkByName(brRtVeth)
		if err != nil {
			if err := ipLinkAddVeth(rtBrVeth, brRtVeth); err != nil {
				return err
			}
			if err := ipLinkSetNetNS(rtBrVeth, netnsName); err != nil {
				return err
			}
		}

		if _, err := netlink.LinkByName(brRtVeth); err == nil {
			log.Println("ip link master")
			if err := ipLinkSetMaster(brRtVeth, brName); err != nil {
				return err
			}
		}

		netnsExec(netnsName, []string{
			"ip", "a", "add", nic.IPv4Address, "dev", rtBrVeth,
		})

		n, err := a.client.SystemV0().Network().Get(vr.Group, vr.Namespace, nic.NetworkID)
		if err != nil {
			return err
		}

		err = netnsExec(netnsName, []string{
			"iptables",
			"-t", "nat",
			"-C", "POSTROUTING",
			"-s", n.Spec.IPv4CIDR,
			"-j", "SNAT",
			"--to-source", natGatewayIP,
		})
		if err != nil {
			err = netnsExec(netnsName, []string{
				"iptables",
				"-t", "nat",
				"-A", "POSTROUTING",
				"-s", n.Spec.IPv4CIDR,
				"-j", "SNAT",
				"--to-source", natGatewayIP,
			})
			log.Println(err)
		}
		netnsExec(netnsName, []string{
			"sh", "-c", `"echo 1 > /proc/sys/net/ipv4/ip_forward"`,
		})
		err = netnsExec(netnsName, []string{
			"ip", "route", "add", "default", "via", vr.Spec.ExternalGateway,
		})
		log.Println(err)
	}

	for _, rule := range vr.Spec.DNATRules {
		log.Println("dnat rule")
		err := netnsExec(netnsName, []string{
			"iptables",
			"-t", "nat",
			"-C", "PREROUTING",
			"-p", "tcp",
			"-i", rtExVeth,
			"-d", natGatewayIP,
			"--dport", fmt.Sprintf("%d", rule.DestPort),
			"-j", "DNAT",
			"--to-destination", fmt.Sprintf("%s:%d", rule.ToDestAddress, rule.ToDestPort),
		})
		if err != nil {
			err := netnsExec(netnsName, []string{
				"iptables",
				"-t", "nat",
				"-A", "PREROUTING",
				"-p", "tcp",
				"-i", rtExVeth,
				"-d", natGatewayIP,
				"--dport", fmt.Sprintf("%d", rule.DestPort),
				"-j", "DNAT",
				"--to-destination", fmt.Sprintf("%s:%d", rule.ToDestAddress, rule.ToDestPort),
			})
			log.Println(err)
		}
	}
	return nil
}

func netnsIsExists(name string) bool {
	_, err := os.Stat(filepath.Join("/var/run/netns", name))
	return err == nil
}

func ipLinkAddVeth(name, peer string) error {
	cmd := exec.Command("ip", "link", "add", name, "type", "veth", "peer", "name", peer)
	return cmd.Run()
}

func ipLinkDel(name string) error {
	cmd := exec.Command("ip", "link", "del", name)
	return cmd.Run()
}

func ipLinkSetNetNS(linkName, netnsName string) error {
	cmd := exec.Command("ip", "link", "set", "up", linkName, "netns", netnsName)
	return cmd.Run()
}

func ipLinkSetMaster(linkName, brName string) error {
	log.Println("ip link set up " + linkName + " master " + brName)
	cmd := exec.Command("ip", "link", "set", "up", linkName, "master", brName)
	return cmd.Run()
}

func netnsAdd(name string) error {
	cmd := exec.Command("ip", "netns", "add", name)
	return cmd.Run()
}

func netnsDel(name string) error {
	cmd := exec.Command("ip", "netns", "del", name)
	return cmd.Run()
}

func netnsExec(name string, command []string) error {
	command = append([]string{
		"netns", "exec", name,
	}, command...)
	cmd := exec.Command("ip", command...)
	return cmd.Run()
}
