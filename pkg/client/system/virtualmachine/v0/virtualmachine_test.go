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
)

const (
	ImageURL = "http://192.168.20.2:8082/bionic-server-cloudimg-amd64.img"
)

func TestVirtualMachineCreate(t *testing.T) {
	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			Name: "test-ns",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	bsClient := bsv0.NewBlockStorageClient("http", "localhost", 8080)
	bs, err := bsClient.Create(&system.BlockStorage{
		Meta: meta.Meta{
			Name:      "test-bs",
			Namespace: ns.ID,
			Annotations: map[string]string{
				"blockstoragev0/type":      "Local",
				"blockstoragev0/node_name": "developvbox",
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
	client := NewVirtualMachineClient("http", "localhost", 8080)
	vm, err := client.Create(&system.VirtualMachine{
		Meta: meta.Meta{
			Name:      "test-vm",
			Namespace: ns.ID,
			Annotations: map[string]string{
				"virtualmachinev0/node_name": "developvbox",
			},
		},
		Spec: system.VirtualMachineSpec{
			RequestMemory: "1G",
			LimitMemory:   "1G",
			RequestVcpus:  "2000m",
			LimitVcpus:    "2000m",
			BlockStorageNames: []string{
				bs.ID,
			},
			ActionState: system.VirtualMachineActionStateStart,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(vm, "", "  ")
	log.Println(string(buf))
}
