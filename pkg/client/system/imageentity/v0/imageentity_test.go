package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

const (
	groupID = "test-group-00"
	imageID = "test-image-00"
)

func TestImageEntityCreate(t *testing.T) {
	//grClient := grv0.NewGroupClient("http", "localhost", 8080)
	//_, err := grClient.Create(&core.Group{
	//	Meta: meta.Meta{
	//		ID:   groupID,
	//		Name: "test-gr",
	//	},
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	client := NewImageEntityClient("http", "localhost", 8080)

	net, err := client.Create(&system.ImageEntity{
		Meta: meta.Meta{
			ID:          imageID,
			Name:        "test-image",
			Group:       groupID,
			Annotations: map[string]string{},
		},
		Spec: system.ImageEntitySpec{
			Hash: "hogehoge",
			Source: system.ImageEntitySource{
				Namespace:      "testns",
				BlockStorageID: "testbs",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(net)
	log.Println(string(buf))

}

func TestImageEntityList(t *testing.T) {
	client := NewImageEntityClient("http", "localhost", 8080)

	netList, err := client.List(groupID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(netList, "", "  ")
	fmt.Println(string(buf))

}

func TestImageEntityGet(t *testing.T) {
	client := NewImageEntityClient("http", "localhost", 8080)

	net, err := client.Get(groupID, imageID)
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.MarshalIndent(net, "", "  ")
	fmt.Println(string(buf))

}

func TestImageEntityUpdate(t *testing.T) {
	client := NewImageEntityClient("http", "localhost", 8080)

	net, err := client.Update(&system.ImageEntity{
		Meta: meta.Meta{
			ID:    imageID,
			Name:  "test-image-changed1",
			Group: groupID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, _ := json.Marshal(net)
	log.Println(string(buf))

}

func TestImageEntityDelete(t *testing.T) {
	client := NewImageEntityClient("http", "localhost", 8080)

	err := client.Delete(groupID, imageID)
	if err != nil {
		t.Fatal(err)
	}

}
