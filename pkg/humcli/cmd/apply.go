package cmd

import (
	"io"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/client"
	"github.com/ophum/humstack/pkg/humcli/cmd/apply"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use: "apply",
	Run: func(cmd *cobra.Command, args []string) {
		clients := client.NewClients(apiServerAddress, apiServerPort)
		applyFuncMap := map[meta.APIType]func(d *yaml.Decoder, clients *client.Clients, debug bool) error{
			meta.APITypeGroupV0:          apply.ApplyGroup,
			meta.APITypeNamespaceV0:      apply.ApplyNamespace,
			meta.APITypeExternalIPPoolV0: apply.ApplyExternalIPPool,
			meta.APITypeExternalIPV0:     apply.ApplyExternalIP,
			meta.APITypeBlockStorageV0:   apply.ApplyBlockStorage,
			meta.APITypeImageV0:          apply.ApplyImage,
			meta.APITypeImageEntityV0:    apply.ApplyImageEntity,
			meta.APITypeNetworkV0:        apply.ApplyNetwork,
			meta.APITypeVirtualMachineV0: apply.ApplyVirtualMachine,
			meta.APITypeVirtualRouterV0:  apply.ApplyVirtualRouter,
		}

		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer f.Close()

			decode := yaml.NewDecoder(f)

			var item meta.Object
			for decode.Decode(&item) == nil {
				r, w := io.Pipe()
				e := yaml.NewEncoder(w)
				go func() {
					e.Encode(item)
					e.Close()
					w.Close()
				}()

				d := yaml.NewDecoder(r)
				if f, ok := applyFuncMap[item.Meta.APIType]; ok {
					if err := f(d, clients, debug); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	},
}
