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
	eippoolID = "test-eippool-00"
)

func TestExternalIPPoolCreate(t *testing.T) {

	client := NewExternalIPPoolClient("http", "localhost", 8080)

	eippool, err := client.Create(&core.ExternalIPPool{
		Meta: meta.Meta{
			ID:   eippoolID,
			Name: "TEST0",
		},
		Spec: core.ExternalIPPoolSpec{
			IPv4CIDR:       "192.168.10.0/24",
			IPv6CIDR:       "fc00::/64",
			BridgeName:     "exBr",
			DefaultGateway: "192.168.10.254",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(eippool)
	log.Println(string(buf))
}

func TestExternalIPPoolList(t *testing.T) {
	client := NewExternalIPPoolClient("http", "localhost", 8080)

	eippoolList, err := client.List()
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(eippoolList, "", "  ")
	fmt.Println(string(buf))

}

func TestExternalIPPoolGet(t *testing.T) {
	client := NewExternalIPPoolClient("http", "localhost", 8080)

	eippool, err := client.Get(eippoolID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(eippool, "", "  ")
	fmt.Println(string(buf))

}

func TestExternalIPPoolUpdate(t *testing.T) {
	client := NewExternalIPPoolClient("http", "localhost", 8080)

	eippool, err := client.Update(&core.ExternalIPPool{
		Meta: meta.Meta{
			Name: "TEST00-changed",
			ID:   eippoolID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(eippool)
	log.Println(string(buf))

}

func TestExternalIPPoolDelete(t *testing.T) {
	client := NewExternalIPPoolClient("http", "localhost", 8080)

	err := client.Delete(eippoolID)
	if err != nil {
		t.Fatal(err)
	}

}
