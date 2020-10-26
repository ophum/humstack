package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyExternalIP(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	eip := &core.ExternalIP{}
	if err := d.Decode(eip); err != nil {
		return err
	}

	old, err := clients.CoreV0().ExternalIP().Get(eip.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		eip, err = clients.CoreV0().ExternalIP().Create(eip)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/externalip/%s created\n", eip.Group, eip.ID)
	} else {
		eip, err = clients.CoreV0().ExternalIP().Update(eip)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/externalip/%s updated\n", eip.Group, eip.ID)
	}

	if debug {
		printYAML(eip)
	}
	return nil
}
