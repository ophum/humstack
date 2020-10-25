package cmd

import (
	"io"
	"log"
	"os"

	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/client"
	"github.com/pkg/errors"
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
				switch item.Meta.APIType {
				case meta.APITypeGroupV0:
					gr := &core.Group{}
					if err = d.Decode(gr); err != nil {
						log.Fatal(errors.Wrap(err, "decode").Error())
					}

					old, err := clients.CoreV0().Group().Get(gr.ID)
					if err != nil {
						log.Fatal(err.Error())
					}
					if old.ID == "" {
						gr, err = clients.CoreV0().Group().Create(gr)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("corev0/group/%s created\n", gr.ID)
					} else {
						gr, err = clients.CoreV0().Group().Update(gr)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("corev0/group/%s updated\n", gr.ID)
					}
				case meta.APITypeNamespaceV0:
					ns := &core.Namespace{}
					if err := d.Decode(ns); err != nil {
						log.Fatal(errors.Wrap(err, "deocde").Error())
					}

					old, err := clients.CoreV0().Namespace().Get(ns.Group, ns.ID)
					if err != nil {
						log.Fatal(err)
					}

					if old.ID == "" {
						ns, err = clients.CoreV0().Namespace().Create(ns)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("%s/corev0/namespace/%s created\n", ns.Group, ns.ID)
					} else {
						ns, err = clients.CoreV0().Namespace().Update(ns)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("%s/corev0/namespace/%s updated\n", ns.Group, ns.ID)
					}
				}
			}
		}
	},
}
