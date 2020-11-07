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
	getCmd.AddCommand(getNodeNetworkCmd)
}

var getNodeNetworkCmd = &cobra.Command{
	Use: "nodenetwork",
	Aliases: []string{
		"nodenet",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		netList, err := clients.SystemV0().NodeNetwork().List(group, namespace)
		if err != nil {
			log.Fatal(err)
		}

		switch output {
		case "json":
			out, err := json.MarshalIndent(netList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(netList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"ID",
				"Name",
				"IPv4CIDR",
				"IPv6CIDR",
				"Network ID",
				"Status",
				"Node",
			})
			for _, n := range netList {
				table.Append([]string{
					n.ID,
					n.Name,
					n.Spec.IPv4CIDR,
					n.Spec.IPv6CIDR,
					n.Spec.ID,
					string(n.Status.State),
					"",
				})
			}

			table.Render()
		}
	},
}
