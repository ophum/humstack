package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	grv0 "github.com/ophum/humstack/pkg/client/core/group/v0"
)

const (
	groupID     = "test-gr"
	namespaceID = "test-namespace-00"
)

func TestNamespaceCreate(t *testing.T) {
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
	client := NewNamespaceClient("http", "localhost", 8080)

	namespace, err := client.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:    namespaceID,
			Group: gr.ID,
			Name:  "TEST0",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(namespace)
	log.Println(string(buf))

}

func TestNamespaceList(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	namespaceList, err := client.List(groupID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(namespaceList, "", "  ")
	fmt.Println(string(buf))

}

func TestNamespaceGet(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	namespace, err := client.Get(groupID, namespaceID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(namespace, "", "  ")
	fmt.Println(string(buf))

}

func TestNamespaceUpdate(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	namespace, err := client.Update(&core.Namespace{
		Meta: meta.Meta{
			Name: "TEST00-changed",
			ID:   namespaceID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(namespace)
	log.Println(string(buf))

}

func TestNamespaceDelete(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	err := client.Delete(groupID, namespaceID)
	if err != nil {
		t.Fatal(err)
	}

}
