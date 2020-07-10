package v0

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
	bsv0 "github.com/ophum/humstack/pkg/client/system/blockstorage/v0"
	netv0 "github.com/ophum/humstack/pkg/client/system/network/v0"
)

const (
	ImageURL = "http://localhost:8082/focal-server-cloudimg-amd64.img"
)

func TestVirtualMachineCreate(t *testing.T) {
	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:   "test-ns",
			Name: "test-ns",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	bsClient := bsv0.NewBlockStorageClient("http", "localhost", 8080)
	bs, err := bsClient.Create(&system.BlockStorage{
		Meta: meta.Meta{
			ID:        "test-bs",
			Name:      "test-bs",
			Namespace: ns.ID,
			Annotations: map[string]string{
				"blockstoragev0/type":      "Local",
				"blockstoragev0/node_name": "X1Carbon",
			},
		},
		Spec: system.BlockStorageSpec{
			RequestSize: "10G",
			LimitSize:   "10G",
			From: system.BlockStorageFrom{
				Type: system.BlockStorageFromTypeHTTP,
				HTTP: system.BlockStorageFromHTTP{
					URL: ImageURL,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	bs2, err := bsClient.Create(&system.BlockStorage{
		Meta: meta.Meta{
			ID:        "test-bs2",
			Name:      "test-bs2",
			Namespace: ns.ID,
			Annotations: map[string]string{
				"blockstoragev0/type":      "Local",
				"blockstoragev0/node_name": "X1Carbon",
			},
		},
		Spec: system.BlockStorageSpec{
			RequestSize: "10G",
			LimitSize:   "10G",
			From: system.BlockStorageFrom{
				Type: system.BlockStorageFromTypeEmpty,
			},
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
			Annotations: map[string]string{
				"networkv0/network_type":    "Bridge",
				"networkv0/default_gateway": "10.0.0.254/24",
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
	client := NewVirtualMachineClient("http", "localhost", 8080)
	vm, err := client.Create(&system.VirtualMachine{
		Meta: meta.Meta{
			ID:        "test-vm",
			Name:      "test-vm",
			Namespace: ns.ID,
			Annotations: map[string]string{
				"virtualmachinev0/node_name": "X1Carbon",
			},
		},
		Spec: system.VirtualMachineSpec{
			RequestMemory: "1G",
			LimitMemory:   "1G",
			RequestVcpus:  "2000m",
			LimitVcpus:    "2000m",
			BlockStorageIDs: []string{
				bs.ID,
				bs2.ID,
			},
			NICs: []*system.VirtualMachineNIC{
				{
					NetworkID:      net.ID,
					IPv4Address:    "10.10.0.1",
					DefaultGateway: "10.10.0.254",
				},
			},
			ActionState: system.VirtualMachineActionStatePowerOn,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(vm, "", "  ")
	log.Println(string(buf))
}

func TestVirtualMachinePowerOff(t *testing.T) {
	client := NewVirtualMachineClient("http", "localhost", 8080)

	vm, err := client.Get("test-ns", "test-vm")
	if err != nil {
		t.Fatal(err)
	}

	vm.Spec.ActionState = system.VirtualMachineActionStatePowerOff
	newVM, err := client.Update(vm)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(newVM, "", "  ")
	log.Println(string(buf))
}

func TestVirtualMachinePowerOn(t *testing.T) {
	client := NewVirtualMachineClient("http", "localhost", 8080)

	vm, err := client.Get("test-ns", "test-vm")
	if err != nil {
		t.Fatal(err)
	}

	vm.Spec.ActionState = system.VirtualMachineActionStatePowerOn
	newVM, err := client.Update(vm)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(newVM, "", "  ")
	log.Println(string(buf))
}
