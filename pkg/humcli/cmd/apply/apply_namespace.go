package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyNamespace(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	ns := &core.Namespace{}
	if err := d.Decode(ns); err != nil {
		return err
	}

	old, err := clients.CoreV0().Namespace().Get(ns.Group, ns.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		ns, err = clients.CoreV0().Namespace().Create(ns)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/namespace/%s created\n", ns.Group, ns.ID)
	} else {
		ns, err = clients.CoreV0().Namespace().Update(ns)
		if err != nil {
			return err
		}
		log.Printf("%s/corev0/namespace/%s updated\n", ns.Group, ns.ID)
	}

	if debug {
		printYAML(ns)
	}
	return nil
}
