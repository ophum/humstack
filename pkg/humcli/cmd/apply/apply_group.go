package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ApplyGroup(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	gr := &core.Group{}
	if err := d.Decode(gr); err != nil {
		log.Fatal(errors.Wrap(err, "decode").Error())
	}

	old, err := clients.CoreV0().Group().Get(gr.ID)
	if err != nil {
		return err
	}
	if old.ID == "" {
		gr, err = clients.CoreV0().Group().Create(gr)
		if err != nil {
			return err
		}
		log.Printf("corev0/group/%s created\n", gr.ID)
	} else {
		gr, err = clients.CoreV0().Group().Update(gr)
		if err != nil {
			return err
		}
		log.Printf("corev0/group/%s updated\n", gr.ID)
	}

	if debug {
		printYAML(gr)
	}

	return nil
}
