package cmd

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use: "create",
}

func init() {
	createCmd.AddCommand(
		diskCreateCmd,
	)
}
