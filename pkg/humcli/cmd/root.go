package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "humstack",
		Short: "humstack cli",
	}
	apiServerAddress string
	apiServerPort    int32
	group            string
	namespace        string
	debug            bool
	output           string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&group, "group", "g", "default", "group id")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "namespace id")
	rootCmd.PersistentFlags().StringVar(&apiServerAddress, "api-server-address", "localhost", "apiserver address")
	rootCmd.PersistentFlags().Int32Var(&apiServerPort, "api-server-port", 8080, "apiserver Port")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format, `table` or `json` or `yaml`")

}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {

	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
