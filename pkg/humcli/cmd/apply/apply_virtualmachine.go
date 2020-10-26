package apply

import (
	"log"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"gopkg.in/yaml.v2"
)

func ApplyVirtualMachine(d *yaml.Decoder, clients *client.Clients, debug bool) error {
	vm := &system.VirtualMachine{}
	if err := d.Decode(vm); err != nil {
		return err
	}

	old, err := clients.SystemV0().VirtualMachine().Get(vm.Group, vm.Namespace, vm.ID)
	if err != nil {
		return err
	}

	if old.ID == "" {
		vm, err = clients.SystemV0().VirtualMachine().Create(vm)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/virtualmachine/%s created\n", vm.Group, vm.Namespace, vm.ID)
	} else {
		vm, err = clients.SystemV0().VirtualMachine().Update(vm)
		if err != nil {
			return err
		}
		log.Printf("%s/%s/systemv0/virtualmachine/%s updated\n", vm.Group, vm.Namespace, vm.ID)
	}

	if debug {
		printYAML(vm)
	}
	return nil
}
