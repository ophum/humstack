package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
)

const (
	namespaceID = "test-namespace-00"
)

func TestNamespaceCreate(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	namespace, err := client.Create(&core.Namespace{
		Meta: meta.Meta{
			ID:   namespaceID,
			Name: "TEST0",
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

	namespaceList, err := client.List()
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(namespaceList, "", "  ")
	fmt.Println(string(buf))

}

func TestNamespaceGet(t *testing.T) {
	client := NewNamespaceClient("http", "localhost", 8080)

	namespace, err := client.Get(namespaceID)
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

	err := client.Delete(namespaceID)
	if err != nil {
		t.Fatal(err)
	}

}
