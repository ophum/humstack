package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyImage(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	image := &system.Image{}
	if err := d.Decode(image); err != nil {
		return err
	}

	old, err := clients.SystemV0().Image().Get(image.Group, image.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		image, err = clients.SystemV0().Image().Create(image)
		if err != nil {
			return err
		}
		log.Printf("%s/systemv0/image/%s created\n", image.Group, image.ID)
	} else {
		image, err = clients.SystemV0().Image().Update(image)
		if err != nil {
			return err
		}
		log.Printf("%s/systemv0/image/%s updated\n", image.Group, image.ID)
	}

	if debug {
		printYAML(image)
	}
	return nil
}
