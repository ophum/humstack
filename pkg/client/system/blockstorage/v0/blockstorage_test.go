package v0

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	grv0 "github.com/ophum/humstack/pkg/client/core/group/v0"
	nsv0 "github.com/ophum/humstack/pkg/client/core/namespace/v0"
)

const (
	imageURL               = "http://192.168.20.2:8082/bionic-server-cloudimg-amd64.img"
	groupID                = "test-group-01"
	groupFromHTTPID        = "test-group-02"
	namespaceID            = "test-namespace-01"
	namespaceFromHTTPID    = "test-namespace-02"
	blockStorageID         = "test-blockstorage-00"
	blockStorageFromHTTPID = "test-blockstorage-01"
)

func TestBlockStorageCreateEmpty(t *testing.T) {

	grClient := grv0.NewGroupClient("http", "localhost", 8080)
	gr, err := grClient.Create(&core.Group{
		Meta: meta.Meta{
			ID:   groupID,
			Name: "test-group",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	ns, err := nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:    namespaceID,
			Group: gr.ID,
			Name:  "test-ns",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client := NewBlockStorageClient("http", "localhost", 8080)
	bs, err := client.Create(&system.BlockStorage{
		Meta: meta.Meta{
			ID:        blockStorageID,
			Name:      "test-bs",
			Namespace: ns.ID,
			Group:     gr.ID,
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

	buf, _ := json.Marshal(bs)
	log.Println(string(buf))
}

func TestBlockStorageCreateHTTP(t *testing.T) {

	grClient := grv0.NewGroupClient("http", "localhost", 8080)
	gr, err := grClient.Create(&core.Group{
		Meta: meta.Meta{
			ID:   groupFromHTTPID,
			Name: "test-group-from-http",
		},
	})
	nsClient := nsv0.NewNamespaceClient("http", "localhost", 8080)
	_, err = nsClient.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:    namespaceFromHTTPID,
			Name:  "test-ns2",
			Group: gr.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	client := NewBlockStorageClient("http", "localhost", 8080)
	bs, err := client.Create(&system.BlockStorage{
		Meta: meta.Meta{
			ID:        blockStorageFromHTTPID,
			Name:      "test-bs-from-http",
			Namespace: namespaceFromHTTPID,
			Group:     gr.ID,
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
					URL: imageURL,
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

	bsList, err := client.List(groupID, namespaceID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(bsList, "", "  ")
	log.Println(string(buf))
}

func TestBlockStorageGet(t *testing.T) {
	client := NewBlockStorageClient("http", "localhost", 8080)

	bsList, err := client.Get(groupFromHTTPID, namespaceFromHTTPID, blockStorageFromHTTPID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(bsList, "", "  ")
	log.Println(string(buf))
}

func TestBlockStorageDeleteState(t *testing.T) {
	client := NewBlockStorageClient("http", "localhost", 8080)

	err := client.DeleteState(groupID, namespaceID, blockStorageID)
	if err != nil {
		t.Fatal(err)
	}

}
