package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
)

const (
	namespaceID = "test-namespace-00"
	networkID   = "test-network-00"
)

func TestNetworkCreate(t *testing.T) {
	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	_, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:   namespaceID,
			Name: "test-ns",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client := NewNetworkClient("http", "localhost", 8080)

	net, err := client.Create(&system.Network{
		Meta: meta.Meta{
			ID:        networkID,
			Name:      "test-network",
			Namespace: namespaceID,
			Annotations: map[string]string{
				"networkv0/network_type":    "Bridge",
				"networkv0/default_gateway": "10.0.0.1/24",
			},
		},
		Spec: system.NetworkSpec{
			IPv4CIDR: "10.0.0.0/24",
			ID:       "100",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(net)
	log.Println(string(buf))

}

func TestNetworkList(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	netList, err := client.List(namespaceID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(netList, "", "  ")
	fmt.Println(string(buf))

}

func TestNetworkGet(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	net, err := client.Get(namespaceID, networkID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(net, "", "  ")
	fmt.Println(string(buf))

}

func TestNetworkUpdate(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	net, err := client.Update(&system.Network{
		Meta: meta.Meta{
			ID:        networkID,
			Name:      "test-network-changed1",
			Namespace: namespaceID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(net)
	log.Println(string(buf))

}

func TestNetworkDelete(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	err := client.Delete(namespaceID, networkID)
	if err != nil {
		t.Fatal(err)
	}

}
