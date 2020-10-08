package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
	agentvmv0 "github.com/ophum/humstack/pkg/agents/system/virtualmachine"
)

func init() {
	getCmd.AddCommand(getVirtualMachineCmd)
}

var getVirtualMachineCmd = &cobra.Command{
	Use: "virtualmachine",
	Aliases: []string{
		"vm",
		"vmachine",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		vmList, err := clients.SystemV0().VirtualMachine().List(group, namespace)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Status(Action=>State)",
			"Vcpus(Limit/Req)\nMemory(Limit/Req)",
			"Node",
			"UUID",
		})
		for _, vm := range vmList {
			table.Append([]string{
				vm.ID,
				vm.Name,
				fmt.Sprintf("%s => %s", vm.Spec.ActionState, vm.Status.State),
				fmt.Sprintf("%s/%s\n%s/%s",
					vm.Spec.LimitVcpus,
					vm.Spec.RequestVcpus,
					vm.Spec.LimitMemory,
					vm.Spec.RequestMemory),
				vm.Annotations[agentvmv0.VirtualMachineV0AnnotationNodeName],
				vm.Spec.UUID,
			})
		}

		table.Render()
	},
}
