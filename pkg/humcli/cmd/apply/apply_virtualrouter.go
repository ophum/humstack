package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyVirtualRouter(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	vr := &system.VirtualRouter{}
	if err := d.Decode(vr); err != nil {
		return err
	}

	old, err := clients.SystemV0().VirtualRouter().Get(vr.Group, vr.Namespace, vr.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		vr, err = clients.SystemV0().VirtualRouter().Create(vr)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/virtualrouter/%s created\n", vr.Group, vr.Namespace, vr.ID)
	} else {
		vr, err = clients.SystemV0().VirtualRouter().Update(vr)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/virtualrouter/%s updated\n", vr.Group, vr.Namespace, vr.ID)
	}

	if debug {
		printYAML(vr)
	}
	return nil
}
