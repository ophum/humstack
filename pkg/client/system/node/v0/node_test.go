package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

func TestNodeCreate(t *testing.T) {
	client := NewNodeClient("http", "localhost", 8080)

	node, err := client.Create(&system.Node{
		Meta: meta.Meta{
			Name: "test-node",
		},
		Spec: system.NodeSpec{
			LimitMemory: "8Gi",
			LimitVcpus:  "10000m",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(node)
	log.Println(string(buf))

}

func TestNodeList(t *testing.T) {
	client := NewNodeClient("http", "localhost", 8080)

	nodeList, err := client.List()
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(nodeList, "", "  ")
	fmt.Println(string(buf))

}

func TestNodeGet(t *testing.T) {
	client := NewNodeClient("http", "localhost", 8080)

	node, err := client.Get("6839fa9d-f6dd-45c9-8e89-6d8f08565dec")
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(node, "", "  ")
	fmt.Println(string(buf))

}

func TestNodeUpdate(t *testing.T) {
	client := NewNodeClient("http", "localhost", 8080)

	node, err := client.Update(&system.Node{
		Meta: meta.Meta{
			ID:   "37825cc2-10b6-4208-9f8d-ec11ee75f9b3",
			Name: "test-node-changed1",
		},
		Spec: system.NodeSpec{
			LimitMemory: "16Gi",
			LimitVcpus:  "10000m",
		},
		Status: system.NodeStatus{
			State: system.NodeStateNotReady,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(node)
	log.Println(string(buf))

}

func TestNodeDelete(t *testing.T) {
	client := NewNodeClient("http", "localhost", 8080)

	err := client.Delete("37825cc2-10b6-4208-9f8d-ec11ee75f9b3")
	if err != nil {
		t.Fatal(err)
	}

}
