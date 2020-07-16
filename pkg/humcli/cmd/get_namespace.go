package cmd

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
)

func init() {
	getCmd.AddCommand(getNamespaceCmd)
}

var getNamespaceCmd = &cobra.Command{
	Use: "namespace",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		nsList, err := clients.CoreV0().Namespace().List(group)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
		})
		for _, n := range nsList {
			table.Append([]string{
				n.ID,
				n.Name,
			})
		}

		table.Render()
	},
}
