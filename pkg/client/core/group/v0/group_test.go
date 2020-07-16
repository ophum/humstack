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
	groupID = "test-group-00"
)

func TestGroupCreate(t *testing.T) {
	client := NewGroupClient("http", "localhost", 8080)

	group, err := client.Create(&core.Group{
		Meta: meta.Meta{
			ID:   groupID,
			Name: "TEST0",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(group)
	log.Println(string(buf))

}

func TestGroupList(t *testing.T) {
	client := NewGroupClient("http", "localhost", 8080)

	groupList, err := client.List()
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(groupList, "", "  ")
	fmt.Println(string(buf))

}

func TestGroupGet(t *testing.T) {
	client := NewGroupClient("http", "localhost", 8080)

	group, err := client.Get(groupID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(group, "", "  ")
	fmt.Println(string(buf))

}

func TestGroupUpdate(t *testing.T) {
	client := NewGroupClient("http", "localhost", 8080)

	group, err := client.Update(&core.Group{
		Meta: meta.Meta{
			Name: "TEST00-changed",
			ID:   groupID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(group)
	log.Println(string(buf))

}

func TestGroupDelete(t *testing.T) {
	client := NewGroupClient("http", "localhost", 8080)

	err := client.Delete(groupID)
	if err != nil {
		t.Fatal(err)
	}

}
