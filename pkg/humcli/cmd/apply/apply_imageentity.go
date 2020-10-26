package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyImageEntity(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	imageEntity := &system.ImageEntity{}
	if err := d.Decode(imageEntity); err != nil {
		return err
	}

	old, err := clients.SystemV0().ImageEntity().Get(imageEntity.Group, imageEntity.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		imageEntity, err = clients.SystemV0().ImageEntity().Create(imageEntity)
		if err != nil {
			return err
		}
		log.Printf("%s/systemv0/imageentity/%s created\n", imageEntity.Group, imageEntity.ID)
	} else {
		imageEntity, err = clients.SystemV0().ImageEntity().Update(imageEntity)
		if err != nil {
			return err
		}
		log.Printf("%s/systemv0/imageentity/%s updated\n", imageEntity.Group, imageEntity.ID)
	}

	if debug {
		printYAML(imageEntity)
	}
	return nil
}
