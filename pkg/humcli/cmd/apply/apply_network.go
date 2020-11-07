package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyNetwork(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	net := &core.Network{}
	if err := d.Decode(net); err != nil {
		return err
	}

	old, err := clients.CoreV0().Network().Get(net.Group, net.Namespace, net.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		net, err = clients.CoreV0().Network().Create(net)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/corev0/network/%s created\n", net.Group, net.Namespace, net.ID)
	} else {
		net, err = clients.CoreV0().Network().Update(net)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/corev0/network/%s updated\n", net.Group, net.Namespace, net.ID)
	}

	if debug {
		printYAML(net)
	}
	return nil
}
