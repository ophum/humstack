package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"
	agentvrv0 "github.com/ophum/humstack/pkg/agents/system/virtualrouter"
)

func init() {
	getCmd.AddCommand(getVirtualRouterCmd)
}

var getVirtualRouterCmd = &cobra.Command{
	Use: "virtualrouter",
	Aliases: []string{
		"vr",
		"vrouter",
	},

	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		vrList, err := clients.SystemV0().VirtualRouter().List(group, namespace)
		if err != nil {
			log.Fatal(err)
		}

		switch output {
		case "json":
			out, err := json.MarshalIndent(vrList, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		case "yaml":
			out, err := yaml.Marshal(vrList)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		default:
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"ID",
				"Name",
				"External Gateway",
				"NAT Gateway IP",
				"EIP => Local",
				"Node",
			})
			for _, vr := range vrList {
				eips := []string{}
				for _, eip := range vr.Spec.ExternalIPs {
					e, err := clients.CoreV0().ExternalIP().Get(eip.ExternalIPID)
					if err != nil {
						log.Fatal(err)
					}

					eips = append(eips, fmt.Sprintf("%s => %s", e.Spec.IPv4Address, eip.BindInternalIPv4Address))
				}
				table.Append([]string{
					vr.ID,
					vr.Name,
					vr.Spec.ExternalGateway,
					vr.Spec.NATGatewayIP,
					strings.Join(eips, "\n"),
					vr.Annotations[agentvrv0.VirtualRouterV0AnnotationNodeName],
				})
			}

			table.Render()
		}
	},
}
