package v0

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	eipv0 "github.com/ophum/humstack/pkg/client/core/externalip/v0"
	eippoolv0 "github.com/ophum/humstack/pkg/client/core/externalippool/v0"
	grv0 "github.com/ophum/humstack/pkg/client/core/group/v0"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
	netv0 "github.com/ophum/humstack/pkg/client/system/network/v0"
)

const (
	imageURL = "http://localhost:8082/focal-server-cloudimg-amd64.img"
	nodeName = "X1Carbon"
)

func TestVirtualRouterCreate(t *testing.T) {
	grClient := grv0.NewGroupClient("http", "localhost", 8080)
	gr, err := grClient.Create(&core.Group{
		Meta: meta.Meta{
			ID:   "test-gr",
			Name: "test-gr",
		},
	})
	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:    "test-ns",
			Name:  "test-ns",
			Group: gr.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	netClient := netv0.NewNetworkClient("http", "localhost", 8080)
	net, err := netClient.Create(&system.Network{
		Meta: meta.Meta{
			ID:        "test-net",
			Name:      "test-net",
			Namespace: ns.ID,
			Group:     gr.ID,
			Annotations: map[string]string{
				"networkv0/network_type": "Bridge",
			},
		},
		Spec: system.NetworkSpec{
			ID:       "100",
			IPv4CIDR: "10.0.0.0/24",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	eippoolClient := eippoolv0.NewExternalIPPoolClient("http", "localhost", 8080)
	eippool, err := eippoolClient.Create(&core.ExternalIPPool{
		Meta: meta.Meta{
			ID:   "testeippool",
			Name: "test eippool",
		},
		Spec: core.ExternalIPPoolSpec{
			IPv4CIDR:       "192.168.10.0/24",
			IPv6CIDR:       "fc00::/64",
			BridgeName:     "exBr",
			DefaultGateway: "192.168.10.254",
		},
	})
	eipClient := eipv0.NewExternalIPClient("http", "localhost", 8080)
	eip, err := eipClient.Create(&core.ExternalIP{
		Meta: meta.Meta{
			ID:   "eip",
			Name: "eip",
		},
		Spec: core.ExternalIPSpec{
			PoolID:      eippool.ID,
			IPv4Address: "192.168.10.1",
			IPv4Prefix:  24,
		},
	})

	client := NewVirtualRouterClient("http", "localhost", 8080)
	vr, err := client.Create(&system.VirtualRouter{
		Meta: meta.Meta{
			ID:        "test-vr",
			Name:      "test-vr",
			Namespace: ns.ID,
			Group:     gr.ID,
			Annotations: map[string]string{
				"virtualrouterv0/node_name": nodeName,
			},
		},
		Spec: system.VirtualRouterSpec{
			ExternalGateway: "192.168.10.254",
			ExternalIPs: []system.VirtualRouterExternalIP{
				{
					ExternalIPID:            eip.ID,
					BindInternalIPv4Address: "10.0.0.1",
				},
			},
			NATGatewayIP: "192.168.10.100/24",
			NICs: []system.VirtualRouterNIC{
				{
					NetworkID:   net.ID,
					IPv4Address: "10.0.0.254/24",
				},
			},
			NATRules: []system.NATRule{
				{
					Type:       system.NATRuleTypeNAPT,
					SrcNetwork: "10.0.0.0/24",
				},
			},
			DNATRules: []system.DNATRule{
				{
					DestPort:      10022,
					ToDestAddress: "10.0.0.2",
					ToDestPort:    22,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(vr, "", "  ")
	log.Println(string(buf))
}

func TestVirtualRouterDeleteState(t *testing.T) {
	client := NewVirtualRouterClient("http", "localhost", 8080)

	err := client.DeleteState("test-gr", "test-ns", "test-vr")
	if err != nil {
		t.Fatal(err)
	}

}
