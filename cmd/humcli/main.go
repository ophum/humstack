package main

import (
	"os"

	"github.com/ophum/humstack/v1/pkg/humcli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
