package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyExternalIPPool(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	eippool := &core.ExternalIPPool{}
	if err := d.Decode(eippool); err != nil {
		return err
	}

	old, err := clients.CoreV0().ExternalIPPool().Get(eippool.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		eippool, err = clients.CoreV0().ExternalIPPool().Create(eippool)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/externalippool/%s created\n", eippool.Group, eippool.ID)
	} else {
		eippool, err = clients.CoreV0().ExternalIPPool().Update(eippool)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/externalippool/%s updated\n", eippool.Group, eippool.ID)
	}

	if debug {
		printYAML(eippool)
	}
	return nil
}
