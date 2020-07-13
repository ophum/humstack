package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use: "create",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		for _, file := range args {
			fmt.Printf("------ CREATE %s ------\n", file)
			data, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			baseData := base{}
			err = yaml.Unmarshal(data, &baseData)
			if err != nil {
				log.Fatal(err)
			}

			switch baseData.APIType {
			case meta.APITypeNamespaceV0:
				ns := &core.Namespace{}
				err = yaml.Unmarshal(data, ns)
				if err != nil {
					log.Fatal(err)
				}

				ns, err = clients.CoreV0().Namespace().Create(ns)
				if err != nil {
					log.Fatal(err)
				}

				nsYAML, err := yaml.Marshal(ns)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(nsYAML))
			case meta.APITypeNetworkV0:
				net := &system.Network{}
				err = yaml.Unmarshal(data, net)
				if err != nil {
					log.Fatal(err)
				}

				net, err = clients.SystemV0().Network().Create(net)
				if err != nil {
					log.Fatal(err)
				}

				netYAML, err := yaml.Marshal(net)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(netYAML))
			case meta.APITypeBlockStorageV0:
				bs := &system.BlockStorage{}
				err = yaml.Unmarshal(data, bs)
				if err != nil {
					log.Fatal(err)
				}

				bs, err = clients.SystemV0().BlockStorage().Create(bs)
				if err != nil {
					log.Fatal(err)
				}

				bsYAML, err := yaml.Marshal(bs)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(bsYAML))
			case meta.APITypeVirtualMachineV0:
				vm := &system.VirtualMachine{}
				err = yaml.Unmarshal(data, vm)
				if err != nil {
					log.Fatal(err)
				}

				vm, err = clients.SystemV0().VirtualMachine().Create(vm)
				if err != nil {
					log.Fatal(err)
				}

				vmYAML, err := yaml.Marshal(vm)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(vmYAML))
			}

		}
	},
}

type base struct {
	meta.Meta `json:"meta" yaml:"meta"`
}
