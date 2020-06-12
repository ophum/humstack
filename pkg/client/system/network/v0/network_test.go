package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

const (
	namespace = "ce6954a8-a10c-4f80-a266-af52c991a968"
	networkID = "6feb0a5c-730f-4ac8-944d-bab10f828f3e"
)

func TestNetworkCreate(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	node, err := client.Create(&system.Network{
		Meta: meta.Meta{
			Name:      "test-network",
			Namespace: namespace,
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

	buf, _ := json.Marshal(node)
	log.Println(string(buf))

}

func TestNetworkList(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	nodeList, err := client.List(namespace)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(nodeList, "", "  ")
	fmt.Println(string(buf))

}

func TestNetworkGet(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	node, err := client.Get(namespace, "test-network")
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(node, "", "  ")
	fmt.Println(string(buf))

}

func TestNetworkUpdate(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	node, err := client.Update(&system.Network{
		Meta: meta.Meta{
			ID:        networkID,
			Name:      "test-network-changed1",
			Namespace: namespace,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(node)
	log.Println(string(buf))

}

func TestNetworkDelete(t *testing.T) {
	client := NewNetworkClient("http", "localhost", 8080)

	err := client.Delete(namespace, networkID)
	if err != nil {
		t.Fatal(err)
	}

}
