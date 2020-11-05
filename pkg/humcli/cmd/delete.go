package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer f.Close()

			decode := yaml.NewDecoder(f)

			var item meta.Object
			for decode.Decode(&item) == nil {
				fmt.Printf("------ DELETE %s %s %s %s ------\n", item.Meta.APIType, item.Meta.Group, item.Meta.Namespace, item.Meta.ID)

				switch item.Meta.APIType {
				case meta.APITypeGroupV0:
					err = clients.CoreV0().Group().DeleteState(item.Meta.ID)
					if err != nil {
						log.Fatal(errors.Wrap(err, "delete").Error())
					}
				case meta.APITypeNamespaceV0:
					err = clients.CoreV0().Namespace().DeleteState(item.Meta.Group, item.Meta.ID)
					if err != nil {
						log.Fatal(errors.Wrap(err, "create").Error())
					}
				case meta.APITypeExternalIPPoolV0:
					err = clients.CoreV0().ExternalIPPool().Delete(item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeExternalIPV0:
					err = clients.CoreV0().ExternalIP().Delete(item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeNetworkV0:
					err = clients.SystemV0().Network().Delete(item.Meta.Group, item.Meta.Namespace, item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeBlockStorageV0:
					err = clients.SystemV0().BlockStorage().DeleteState(item.Meta.Group, item.Meta.Namespace, item.Meta.ID)
					if err != nil {
						log.Fatal(err.Error())
					}
				case meta.APITypeVirtualMachineV0:
					err = clients.SystemV0().VirtualMachine().DeleteState(item.Meta.Group, item.Meta.Namespace, item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeVirtualRouterV0:
					err = clients.SystemV0().VirtualRouter().DeleteState(item.Meta.Group, item.Meta.Namespace, item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeImageV0:
					err = clients.SystemV0().Image().Delete(item.Meta.Group, item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				case meta.APITypeImageEntityV0:
					err = clients.SystemV0().ImageEntity().Delete(item.Meta.Group, item.Meta.ID)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}

	},
}
