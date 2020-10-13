package cmd

import (
	"log"

	"github.com/ophum/humstack/pkg/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use: "watch",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients("localhost", 8080)

		clients.WatchV0().Watch(func(before, after interface{}) {
			log.Println("WATCH")
			log.Printf("BEFORE: %+v", before)
			log.Printf("AFTER: %+v", after)
		})

	},
}
