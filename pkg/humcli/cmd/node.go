package cmd

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var nodeListCmd = &cobra.Command{
	Use: "nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newNodeClient()
		if err != nil {
			return err
		}
		nodes, err := c.List(cmd.Context())
		if err != nil {
			return errors.Wrap(err, "failed to get nodes")
		}
		fmt.Println("name\thostname\tstatus")
		for _, node := range nodes {

			fmt.Printf("%s\t%s\t%s\n", node.Name, node.Hostname, node.Status)
			if len(node.Agents) > 0 {
				for i, a := range node.Agents {
					if i == len(node.Agents)-1 {
						fmt.Print(" └ ")
					} else {
						fmt.Print(" │ ")
					}
					fmt.Printf("%s: %s %s\n", a.Name, a.Command, strings.Join(a.Args, " "))
				}
			} else {
				fmt.Println("")
			}
		}
		return nil
	},
}
