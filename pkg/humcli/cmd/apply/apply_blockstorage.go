package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyBlockStorage(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	bs := &system.BlockStorage{}
	if err := d.Decode(bs); err != nil {
		return err
	}

	old, err := clients.SystemV0().BlockStorage().Get(bs.Group, bs.Namespace, bs.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		bs, err = clients.SystemV0().BlockStorage().Create(bs)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/blockstorage/%s created\n", bs.Group, bs.Namespace, bs.ID)
	} else {
		bs, err = clients.SystemV0().BlockStorage().Update(bs)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/blockstorage/%s updated\n", bs.Group, bs.Namespace, bs.ID)
	}

	if debug {
		printYAML(bs)
	}
	return nil
}
