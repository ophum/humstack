package cmd

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

func init() {
	getCmd.AddCommand(getNetworkCmd)
}

var getNetworkCmd = &cobra.Command{
	Use: "network",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		netList, err := clients.SystemV0().Network().List(namespace)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"IPv4CIDR",
			"IPv6CIDR",
			"Network ID",
		})
		for _, n := range netList {
			table.Append([]string{
				n.ID,
				n.Name,
				n.Spec.IPv4CIDR,
				n.Spec.IPv6CIDR,
				n.Spec.ID,
			})
		}

		table.Render()
	},
}
