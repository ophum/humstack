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
	getCmd.AddCommand(getExternalIPPoolCmd)
}

var getExternalIPPoolCmd = &cobra.Command{
	Use: "externalippool",
	Aliases: []string{
		"eippool",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		eippoolList, err := clients.CoreV0().ExternalIPPool().List()
		if err != nil {
			log.Fatal(err)
		}

		switch output {
		case "json":
			out, err := json.MarshalIndent(eippoolList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(eippoolList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
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
		}
	},
}
