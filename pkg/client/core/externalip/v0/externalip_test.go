package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	eippoolv0 "github.com/ophum/humstack/pkg/client/core/externalippool/v0"
)

const (
	eipID     = "test-eip-00"
	eippoolID = "test-eippool-00"
)

func TestExternalIPCreate(t *testing.T) {
	eippoolClient := eippoolv0.NewExternalIPPoolClient("http", "localhost", 8080)
	pool, err := eippoolClient.Create(&core.ExternalIPPool{
		Meta: meta.Meta{
			ID:   eippoolID,
			Name: "test pool",
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

	client := NewExternalIPClient("http", "localhost", 8080)

	eip, err := client.Create(&core.ExternalIP{
		Meta: meta.Meta{
			ID:   eipID,
			Name: "TEST0",
		},
		Spec: core.ExternalIPSpec{
			PoolID:      pool.ID,
			IPv4Address: "192.168.10.1",
			IPv4Prefix:  24,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(eip)
	log.Println(string(buf))
}

func TestExternalIPList(t *testing.T) {
	client := NewExternalIPClient("http", "localhost", 8080)

	eipList, err := client.List()
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(eipList, "", "  ")
	fmt.Println(string(buf))

}

func TestExternalIPGet(t *testing.T) {
	client := NewExternalIPClient("http", "localhost", 8080)

	eip, err := client.Get(eipID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(eip, "", "  ")
	fmt.Println(string(buf))

}

func TestExternalIPUpdate(t *testing.T) {
	client := NewExternalIPClient("http", "localhost", 8080)

	eip, err := client.Update(&core.ExternalIP{
		Meta: meta.Meta{
			Name: "TEST00-changed",
			ID:   eipID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(eip)
	log.Println(string(buf))

}

func TestExternalIPDelete(t *testing.T) {
	client := NewExternalIPClient("http", "localhost", 8080)

	err := client.Delete(eipID)
	if err != nil {
		t.Fatal(err)
	}

}
