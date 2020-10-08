package cmd

import (
	"log"
	"os"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"

	"github.com/olekukonko/tablewriter"
	agentbsv0 "github.com/ophum/humstack/pkg/agents/system/blockstorage"
)

func init() {
	getCmd.AddCommand(getBlockStorageCmd)
}

var getBlockStorageCmd = &cobra.Command{
	Use: "blockstorage",
	Aliases: []string{
		"bs",
	},
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)
		bsList, err := clients.SystemV0().BlockStorage().List(group, namespace)
		if err != nil {
			log.Fatal(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"RequestSize",
			"LimitSize",
			"NodeName",
			"Type",
			"FromType",
			"Status",
		})
		for _, bs := range bsList {
			state := string(bs.Status.State)
			if bs.DeleteState != "" {
				state = string(bs.DeleteState)
			}
			table.Append([]string{
				bs.ID,
				bs.Name,
				bs.Spec.RequestSize,
				bs.Spec.LimitSize,
				bs.Annotations[agentbsv0.BlockStorageV0AnnotationNodeName],
				bs.Annotations[agentbsv0.BlockStorageV0AnnotationType],
				string(bs.Spec.From.Type),
				state,
			})
		}

		table.Render()
	},
}
