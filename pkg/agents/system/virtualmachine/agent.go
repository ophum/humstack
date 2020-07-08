package virtualmachine

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
)

type VirtualMachineAgent struct {
	client *client.Clients
}

const (
	VirtualMachineV0AnnotationNodeName = "virtualmachinev0/node_name"
)

func NewVirtualMachineAgent(client *client.Clients) *VirtualMachineAgent {
	return &VirtualMachineAgent{
		client: client,
	}
}

func (a *VirtualMachineAgent) Run() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	nodeName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ticker.C:
			nsList, err := a.client.CoreV0().Namespace().List()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, ns := range nsList {
				vmList, err := a.client.SystemV0().VirtualMachine().List(ns.ID)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, vm := range vmList {
					oldHash := vm.ResourceHash
					if vm.Annotations[VirtualMachineV0AnnotationNodeName] != nodeName {
						continue
					}
					err = a.syncVirtualMachine(vm)
					if err != nil {
						log.Println(err)
						continue
					}

					if vm.ResourceHash == oldHash {
						log.Printf("vm(`%s`) no update.\n", vm.ID)
					}

					_, err := a.client.SystemV0().VirtualMachine().Update(vm)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			}

		}
	}
}

func (a *VirtualMachineAgent) syncVirtualMachine(vm *system.VirtualMachine) error {
	if vm.Status.State == system.VirtualMachineStateRunning {
		return nil
	}

	vcpus := withUnitToWithoutUnit(vm.Spec.LimitVcpus)
	command := "qemu-system-x86_64"
	args := []string{
		"-enable-kvm",
		"-uuid",
		vm.ID,
		"-name",
		fmt.Sprintf("guest=%s,debug-threads=on", vm.Name),
		"-daemonize",
		"-nodefaults",
		"-vnc",
		fmt.Sprintf("0.0.0.0:1"),
		"-smp",
		fmt.Sprintf("%s,sockets=1,cores=%s,threads=1", vcpus, vcpus),
		"-cpu",
		"host",
		"-m",
		vm.Spec.LimitMemory,
		"-device",
		"VGA,id=video0,bus=pci.0",
		filepath.Join("./blockstorages", vm.Namespace, vm.Spec.BlockStorageNames[0]),
	}

	log.Printf("create vm `%s`", vm.ID)
	log.Println(command, args)
	cmd := exec.Command(command, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Println(err.Error())
		return err
	}

	vm.Status.State = system.VirtualMachineStateRunning
	return setHash(vm)
}

func setHash(vm *system.VirtualMachine) error {
	vm.ResourceHash = ""
	resourceJSON, err := json.Marshal(vm)
	if err != nil {
		return err
	}

	hash := md5.Sum(resourceJSON)
	vm.ResourceHash = fmt.Sprintf("%x", hash)
	return nil
}

const (
	UnitGigabyte = 'G'
	UnitMegabyte = 'M'
	UnitKilobyte = 'K'
	UnitMilli    = 'm'
)

func withUnitToWithoutUnit(numberWithUnit string) string {
	length := len(numberWithUnit)
	if numberWithUnit[length-1] >= '0' && numberWithUnit[length-1] <= '9' {
		return numberWithUnit
	}

	number, err := strconv.ParseInt(numberWithUnit[:length-1], 10, 64)
	if err != nil {
		return "0"
	}

	switch numberWithUnit[length-1] {
	case UnitGigabyte:
		return fmt.Sprintf("%d", number*1024*1024*1024)
	case UnitMegabyte:
		return fmt.Sprintf("%d", number*1024*1024)
	case UnitKilobyte:
		return fmt.Sprintf("%d", number*1024)
	case UnitMilli:
		return fmt.Sprintf("%d", number/1000)
	}
	return "0"
}
