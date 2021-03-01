package virtualrouter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ophum/humstack/pkg/agents/system/nodenetwork/utils"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/vishvananda/netlink"
	"go.uber.org/zap"
)

type VirtualRouterAgent struct {
	client *client.Clients
	logger *zap.Logger

	externalBridge  string
	floatingIPCIDR  string
	usedFloatingIPs map[string]bool
}

const (
	VirtualRouterV0AnnotationNodeName = "virtualrouterv0/node_name"
)

func NewVirtualRouterAgent(client *client.Clients, externalBridge string, floatingIPCIDR string, usedFloatingIPs []string, logger *zap.Logger) *VirtualRouterAgent {
	return &VirtualRouterAgent{
		client:          client,
		logger:          logger,
		externalBridge:  externalBridge,
		floatingIPCIDR:  floatingIPCIDR,
		usedFloatingIPs: map[string]bool{},
	}
}

func (a *VirtualRouterAgent) Run(pollingDuration time.Duration) {
	ticker := time.NewTicker(pollingDuration)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		a.logger.Panic(
			"get hostname",
			zap.String("msg", err.Error()),
			zap.Time("time", time.Now()),
		)
	}

	for {
		select {
		case <-ticker.C:
			grList, err := a.client.CoreV0().Group().List()
			if err != nil {
				a.logger.Error(
					"get group list",
					zap.String("msg", err.Error()),
					zap.Time("time", time.Now()),
				)
				continue
			}

			for _, group := range grList {
				nsList, err := a.client.CoreV0().Namespace().List(group.ID)
				if err != nil {
					a.logger.Error(
						"get namespace list",
						zap.String("msg", err.Error()),
						zap.Time("time", time.Now()),
					)
					continue
				}

				for _, ns := range nsList {
					vrList, err := a.client.SystemV0().VirtualRouter().List(group.ID, ns.ID)
					if err != nil {
						a.logger.Error(
							"get virtualrouter list",
							zap.String("msg", err.Error()),
							zap.Time("time", time.Now()),
						)
						continue
					}

					for _, vr := range vrList {
						oldHash := vr.ResourceHash
						if vr.Annotations[VirtualRouterV0AnnotationNodeName] != nodeName {
							continue
						}

						err = a.syncVirtualRouter(vr)
						if err != nil {
							a.logger.Error(
								"sync virtualrouter",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}

						if vr.ResourceHash == oldHash {
							continue
						}

						_, err := a.client.SystemV0().VirtualRouter().Update(vr)
						if err != nil {
							a.logger.Error(
								"update virtualrouter",
								zap.String("msg", err.Error()),
								zap.Time("time", time.Now()),
							)
							continue
						}
					}
				}
			}
		}
	}
}

func (a *VirtualRouterAgent) syncVirtualRouter(vr *system.VirtualRouter) error {

	netnsName := utils.GenerateName("netns-", vr.Group+vr.Namespace+vr.ID)
	rtExVeth := utils.GenerateName("hrt-ex-", netnsName)
	exRtVeth := utils.GenerateName("hex-rt-", netnsName)

	natRule := []string{"*nat"}
	if !netnsIsExists(netnsName) {
		if err := netnsAdd(netnsName); err != nil {
			return err
		}

		if err := ipLinkAddVeth(rtExVeth, exRtVeth); err != nil {
			return err
		}

		if err := ipLinkSetNetNS(rtExVeth, netnsName); err != nil {
			return err
		}

		if err := ipLinkSetMaster(exRtVeth, a.externalBridge); err != nil {
			return err
		}

		netnsExec(netnsName, []string{
			"ip", "a", "add", vr.Spec.NATGatewayIP, "dev", rtExVeth,
		})
	}

	for _, e := range vr.Spec.ExternalIPs {
		eip, err := a.client.CoreV0().ExternalIP().Get(e.ExternalIPID)
		if err != nil {
			return err
		}

		netnsExec(netnsName, []string{
			"ip", "a",
			"add",
			fmt.Sprintf("%s/%d",
				eip.Spec.IPv4Address,
				eip.Spec.IPv4Prefix,
			),
			"dev", rtExVeth,
		})

		// iptables -t nat -A PREROUTING -d ${daddr} -j DNAT --to-destination ${DEST}
		// iptables -t nat -A POSTROUTING -s ${saddr} -j SNAT --to-source ${daddr}

		daddr := eip.Spec.IPv4Address

		natRule = append(natRule, fmt.Sprintf("-A PREROUTING -d %s -j DNAT --to-destination %s", daddr, e.BindInternalIPv4Address))
		natRule = append(natRule, fmt.Sprintf("-A POSTROUTING -s %s -j SNAT --to-source %s", e.BindInternalIPv4Address, daddr))
	}

	natGatewayIP := strings.Split(vr.Spec.NATGatewayIP, "/")[0]
	for _, nic := range vr.Spec.NICs {
		rtBrVeth := utils.GenerateName("hrt-br-", netnsName+nic.NetworkID)
		brRtVeth := utils.GenerateName("hbr-rt-", netnsName+nic.NetworkID)
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
			if err := ipLinkSetMaster(brRtVeth, brName); err != nil {
				return err
			}
		}

		netnsExec(netnsName, []string{
			"ip", "a", "add", nic.IPv4Address, "dev", rtBrVeth,
		})

		n, err := a.client.CoreV0().Network().Get(vr.Group, vr.Namespace, nic.NetworkID)
		if err != nil {
			return err
		}

		natRule = append(natRule, fmt.Sprintf("-A POSTROUTING -s %s -j SNAT --to-source %s", n.Spec.Template.Spec.IPv4CIDR, natGatewayIP))
		netnsExec(netnsName, []string{
			"sh", "-c", `"echo 1 > /proc/sys/net/ipv4/ip_forward"`,
		})
		err = netnsExec(netnsName, []string{
			"ip", "route", "add", "default", "via", vr.Spec.ExternalGateway,
		})
	}

	for _, rule := range vr.Spec.DNATRules {
		natRule = append(natRule,
			fmt.Sprintf("-A PREROUTING -p tcp -i %s -d %s --dport %d -j DNAT --to-destination %s:%d",
				rtExVeth,
				natGatewayIP,
				rule.DestPort,
				rule.ToDestAddress, rule.ToDestPort))
	}

	natRule = append(natRule, "COMMIT")
	natFile := strings.Join(natRule, "\n")
	cmd := exec.Command("ip", "netns", "exec", netnsName, "iptables-restore")
	w, _ := cmd.StdinPipe()
	_, err := io.WriteString(w, natFile+"\n")
	w.Close()
	if err != nil {
		return err
	}

	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	vr.Status.State = system.VirtualRouterStateRunning
	return setHash(vr)
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

func setHash(vr *system.VirtualRouter) error {
	vr.ResourceHash = ""
	resourceJSON, err := json.Marshal(vr)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	vr.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}
