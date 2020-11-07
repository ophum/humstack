package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyNodeNetwork(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	net := &system.NodeNetwork{}
	if err := d.Decode(net); err != nil {
		return err
	}

	old, err := clients.SystemV0().NodeNetwork().Get(net.Group, net.Namespace, net.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		net, err = clients.SystemV0().NodeNetwork().Create(net)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/nodenetwork/%s created\n", net.Group, net.Namespace, net.ID)
	} else {
		net, err = clients.SystemV0().NodeNetwork().Update(net)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/nodenetwork/%s updated\n", net.Group, net.Namespace, net.ID)
	}

	if debug {
		printYAML(net)
	}
	return nil
}
