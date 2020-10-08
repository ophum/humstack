package cmd

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

func init() {
	getCmd.AddCommand(getExternalIPPoolCmd)
}

var getExternalIPPoolCmd = &cobra.Command{
	Use: "externalippool",
	Aliases: []string{
		"eippool",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		eippoolList, err := clients.CoreV0().ExternalIPPool().List()
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"Bridge",
			"IPv4 CIDR",
			"IPv6 CIDR",
		})
		for _, eippool := range eippoolList {
			table.Append([]string{
				eippool.Name,
				eippool.Spec.BridgeName,
				eippool.Spec.IPv4CIDR,
				eippool.Spec.IPv6CIDR,
			})
		}

		table.Render()
	},
}
