package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use: "create",
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
				fmt.Printf("------ CREATE %s %s %s %s ------\n", item.Meta.APIType, item.Meta.Group, item.Meta.Namespace, item.Meta.ID)

				// よくない実装だと思う...
				r, w := io.Pipe()
				e := yaml.NewEncoder(w)
				go func() {
					e.Encode(item)
					e.Close()
					w.Close()
				}()

				d := yaml.NewDecoder(r)
				switch item.Meta.APIType {
				case meta.APITypeGroupV0:
					gr := &core.Group{}
					if err = d.Decode(gr); err != nil {
						log.Fatal(errors.Wrap(err, "decode").Error())
					}

					gr, err = clients.CoreV0().Group().Create(gr)
					if err != nil {
						log.Fatal(errors.Wrap(err, "create").Error())
					}

					resYAML, err := yaml.Marshal(gr)
					if err != nil {
						log.Fatal(errors.Wrap(err, "marshal").Error())
					}
					fmt.Println(string(resYAML))
				case meta.APITypeNamespaceV0:
					ns := &core.Namespace{}
					if err = d.Decode(ns); err != nil {
						log.Fatal(errors.Wrap(err, "decode").Error())
					}

					ns, err = clients.CoreV0().Namespace().Create(ns)
					if err != nil {
						log.Fatal(errors.Wrap(err, "create").Error())
					}

					resYAML, err := yaml.Marshal(ns)
					if err != nil {
						log.Fatal(errors.Wrap(err, "marshal").Error())
					}
					fmt.Println(string(resYAML))

				case meta.APITypeExternalIPPoolV0:
					eippool := &core.ExternalIPPool{}
					err = d.Decode(eippool)
					if err != nil {
						log.Fatal(err)
					}

					eippool, err = clients.CoreV0().ExternalIPPool().Create(eippool)
					if err != nil {
						log.Fatal(err)
					}

					eippoolYAML, err := yaml.Marshal(eippool)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(eippoolYAML))
				case meta.APITypeExternalIPV0:
					eip := &core.ExternalIP{}
					err = d.Decode(eip)
					if err != nil {
						log.Fatal(err)
					}

					eip, err = clients.CoreV0().ExternalIP().Create(eip)
					if err != nil {
						log.Fatal(err)
					}

					eipYAML, err := yaml.Marshal(eip)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(eipYAML))
				case meta.APITypeNetworkV0:
					net := &system.Network{}
					err = d.Decode(net)
					if err != nil {
						log.Fatal(err)
					}

					net, err = clients.SystemV0().Network().Create(net)
					if err != nil {
						log.Fatal(err)
					}

					netYAML, err := yaml.Marshal(net)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(netYAML))
				case meta.APITypeBlockStorageV0:
					bs := &system.BlockStorage{}
					err = d.Decode(bs)
					if err != nil {
						log.Fatal(err)
					}

					log.Println(bs)
					bs, err = clients.SystemV0().BlockStorage().Create(bs)
					if err != nil {
						log.Fatal(err.Error())
					}

					bsYAML, err := yaml.Marshal(bs)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(bsYAML))
				case meta.APITypeVirtualMachineV0:
					vm := &system.VirtualMachine{}
					err = d.Decode(vm)
					if err != nil {
						log.Fatal(err)
					}

					vm, err = clients.SystemV0().VirtualMachine().Create(vm)
					if err != nil {
						log.Fatal(err)
					}

					vmYAML, err := yaml.Marshal(vm)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(vmYAML))
				case meta.APITypeVirtualRouterV0:
					vr := &system.VirtualRouter{}
					err = d.Decode(vr)
					if err != nil {
						log.Fatal(err)
					}

					vr, err = clients.SystemV0().VirtualRouter().Create(vr)
					if err != nil {
						log.Fatal(err)
					}

					vrYAML, err := yaml.Marshal(vr)
					if err != nil {
						log.Fatal(err)
					}

					fmt.Println(string(vrYAML))

				}
				r.Close()
			}
		}

	},
}
