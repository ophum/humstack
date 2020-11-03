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
	getCmd.AddCommand(getNodeCmd)
}

var getNodeCmd = &cobra.Command{
	Use: "node",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		nodeList, err := clients.SystemV0().Node().List()
		if err != nil {
			log.Fatal(err)
		}

		switch output {
		case "json":
			out, err := json.MarshalIndent(nodeList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(nodeList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
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
		}
	},
}
