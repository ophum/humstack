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
	getCmd.AddCommand(getImageCmd)
}

var getImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		imList, err := clients.SystemV0().Image().List(group)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Tags (EntityID)",
		})
		for _, im := range imList {

			tags := ""
			for tag, entityID := range im.Spec.EntityMap {
				tags += fmt.Sprintf("%s (%s)\n", tag, entityID)
			}

			table.Append([]string{
				im.ID,
				im.Name,
				tags,
			})
		}

		table.Render()
	},
}
