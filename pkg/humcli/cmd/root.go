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
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVar(&group, "g", "default", "group id")
	rootCmd.PersistentFlags().StringVar(&namespace, "n", "default", "namespace id")
	rootCmd.PersistentFlags().StringVar(&apiServerAddress, "api-server-address", "localhost", "apiserver address")
	rootCmd.PersistentFlags().Int32Var(&apiServerPort, "api-server-port", 8080, "apiserver Port")

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
