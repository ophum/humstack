package cmd

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

func init() {
	getCmd.AddCommand(getNodeCmd)
}

var getNodeCmd = &cobra.Command{
	Use: "node",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		nodeList, err := clients.SystemV0().Node().List()
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"LimitVcpus",
			"LimitMemory",
		})
		for _, n := range nodeList {
			table.Append([]string{
				n.Name,
				n.Spec.LimitVcpus,
				n.Spec.LimitMemory,
			})
		}

		table.Render()
	},
}
