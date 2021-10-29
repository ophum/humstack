package cmd

import (
	"net/url"

	"github.com/ophum/humstack/v1/pkg/client"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "humcli",
	Short: "humstack cli tool",
}

var (
	apiEndpoint string
)

func newDiskClient() (*client.DiskClient, error) {
	endpoint, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil, err
	}
	c := client.NewDiskClient(*endpoint)
	return c, nil
}
func init() {
	RootCmd.PersistentFlags().StringVarP(&apiEndpoint, "apiserver-endpoint", "H", "http://localhost:8080", "api server endpoint url, ex: http://localhost:8080")
	cobra.OnInitialize()
	RootCmd.AddCommand(
		listCmd,
		createCmd,
	)
}
