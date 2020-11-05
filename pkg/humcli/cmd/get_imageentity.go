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

		switch output {
		case "json":
			out, err := json.MarshalIndent(ieList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(ieList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
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
		}
	},
}
