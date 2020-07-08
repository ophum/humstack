package v0

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
)

const (
	ImageURL = "http://192.168.20.2:8082/bionic-server-cloudimg-amd64.img"
)

func TestBlockStorageCreateEmpty(t *testing.T) {

	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			Name: "test-ns",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client := NewBlockStorageClient("http", "localhost", 8080)
	bs, err := client.Create(&system.BlockStorage{
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
				Type: system.BlockStorageFromTypeEmpty,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(bs)
	log.Println(string(buf))
}

func TestBlockStorageCreateHTTP(t *testing.T) {

	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			Name: "test-ns2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client := NewBlockStorageClient("http", "localhost", 8080)
	bs, err := client.Create(&system.BlockStorage{
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

	buf, _ := json.Marshal(bs)
	log.Println(string(buf))
}
func TestBlockStorageList(t *testing.T) {
	client := NewBlockStorageClient("http", "localhost", 8080)

	bsList, err := client.List("de4932c3-323f-464a-979f-6d037465a0bf")
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(bsList, "", "  ")
	log.Println(string(buf))
}

func TestBlockStorageGet(t *testing.T) {
	client := NewBlockStorageClient("http", "localhost", 8080)

	bsList, err := client.Get("de4932c3-323f-464a-979f-6d037465a0bf", "90976489-70b9-435f-9ac2-61e681822625")
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(bsList, "", "  ")
	log.Println(string(buf))
}
