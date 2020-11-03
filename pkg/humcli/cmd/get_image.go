package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

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

		switch output {
		case "json":
			out, err := json.MarshalIndent(imList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(imList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
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
		}
	},
}
