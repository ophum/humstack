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
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use: "update",
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
				fmt.Printf("------ UPDATE %s %s %s %s ------\n", item.Meta.APIType, item.Meta.Group, item.Meta.Namespace, item.Meta.ID)

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

					gr, err = clients.CoreV0().Group().Update(gr)
					if err != nil {
						log.Fatal(errors.Wrap(err, "create").Error())
					}

					printYAML(gr)
				case meta.APITypeNamespaceV0:
					ns := &core.Namespace{}
					if err = d.Decode(ns); err != nil {
						log.Fatal(errors.Wrap(err, "decode").Error())
					}

					ns, err = clients.CoreV0().Namespace().Update(ns)
					if err != nil {
						log.Fatal(errors.Wrap(err, "create").Error())
					}

					printYAML(ns)
				case meta.APITypeExternalIPPoolV0:
					eippool := &core.ExternalIPPool{}
					err = d.Decode(eippool)
					if err != nil {
						log.Fatal(err)
					}

					eippool, err = clients.CoreV0().ExternalIPPool().Update(eippool)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(eippool)
				case meta.APITypeExternalIPV0:
					eip := &core.ExternalIP{}
					err = d.Decode(eip)
					if err != nil {
						log.Fatal(err)
					}

					eip, err = clients.CoreV0().ExternalIP().Update(eip)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(eip)
				case meta.APITypeNetworkV0:
					net := &system.Network{}
					err = d.Decode(net)
					if err != nil {
						log.Fatal(err)
					}

					net, err = clients.SystemV0().Network().Update(net)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(net)
				case meta.APITypeBlockStorageV0:
					bs := &system.BlockStorage{}
					err = d.Decode(bs)
					if err != nil {
						log.Fatal(err)
					}

					bs, err = clients.SystemV0().BlockStorage().Update(bs)
					if err != nil {
						log.Fatal(err.Error())
					}

					printYAML(bs)
				case meta.APITypeVirtualMachineV0:
					vm := &system.VirtualMachine{}
					err = d.Decode(vm)
					if err != nil {
						log.Fatal(err)
					}

					vm, err = clients.SystemV0().VirtualMachine().Update(vm)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(vm)
				case meta.APITypeVirtualRouterV0:
					vr := &system.VirtualRouter{}
					err = d.Decode(vr)
					if err != nil {
						log.Fatal(err)
					}

					vr, err = clients.SystemV0().VirtualRouter().Update(vr)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(vr)
				case meta.APITypeImageV0:
					image := &system.Image{}
					if err = d.Decode(image); err != nil {
						log.Fatal(err)
					}

					image, err = clients.SystemV0().Image().Update(image)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(image)
				case meta.APITypeImageEntityV0:
					imageEntity := &system.ImageEntity{}
					if err := d.Decode(imageEntity); err != nil {
						log.Fatal(err)
					}

					imageEntity, err = clients.SystemV0().ImageEntity().Update(imageEntity)
					if err != nil {
						log.Fatal(err)
					}

					printYAML(imageEntity)
				}
				r.Close()
			}
		}

	},
}
