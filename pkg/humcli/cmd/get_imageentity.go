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
	getCmd.AddCommand(getImageEntityCmd)
}

var getImageEntityCmd = &cobra.Command{
	Use: "imageentity",
	Aliases: []string{
		"ie",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		ieList, err := clients.SystemV0().ImageEntity().List(group)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Source(ns/bs)",
			"Status",
			"Hash",
		})
		for _, ie := range ieList {
			state := string(ie.Status.State)
			if ie.DeleteState != "" {
				state = string(ie.DeleteState)
			}
			table.Append([]string{
				ie.ID,
				ie.Name,
				fmt.Sprintf("%s/%s",
					ie.Spec.Source.Namespace,
					ie.Spec.Source.BlockStorageID),
				state,
				ie.Spec.Hash,
			})
		}

		table.Render()
	},
}
