package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

func init() {
	getCmd.AddCommand(getExternalIPCmd)
}

var getExternalIPCmd = &cobra.Command{
	Use: "externalip",
	Aliases: []string{
		"eip",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		eipList, err := clients.CoreV0().ExternalIP().List()
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"Name",
			"Pool ID",
			"IPv4",
			"IPv6",
		})
		for _, eip := range eipList {
			table.Append([]string{
				eip.Name,
				eip.Spec.PoolID,
				fmt.Sprintf("%s/%d", eip.Spec.IPv4Address, eip.Spec.IPv4Prefix),
				fmt.Sprintf("%s/%d", eip.Spec.IPv6Address, eip.Spec.IPv6Prefix),
			})
		}

		table.Render()
	},
}
