package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/ophum/humstack/cmd/apiserver/statik"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core/externalip"
	eipv0 "github.com/ophum/humstack/pkg/api/core/externalip/v0"
	"github.com/ophum/humstack/pkg/api/core/externalippool"
	eippoolv0 "github.com/ophum/humstack/pkg/api/core/externalippool/v0"
	"github.com/ophum/humstack/pkg/api/core/group"
	grv0 "github.com/ophum/humstack/pkg/api/core/group/v0"
	"github.com/ophum/humstack/pkg/api/core/namespace"
	nsv0 "github.com/ophum/humstack/pkg/api/core/namespace/v0"
	"github.com/ophum/humstack/pkg/api/core/network"
	netv0 "github.com/ophum/humstack/pkg/api/core/network/v0"
	"github.com/ophum/humstack/pkg/api/system/blockstorage"
	bsv0 "github.com/ophum/humstack/pkg/api/system/blockstorage/v0"
	"github.com/ophum/humstack/pkg/api/system/image"
	imv0 "github.com/ophum/humstack/pkg/api/system/image/v0"
	"github.com/ophum/humstack/pkg/api/system/imageentity"
	iev0 "github.com/ophum/humstack/pkg/api/system/imageentity/v0"
	"github.com/ophum/humstack/pkg/api/system/node"
	nodev0 "github.com/ophum/humstack/pkg/api/system/node/v0"
	"github.com/ophum/humstack/pkg/api/system/nodenetwork"
	nodenetv0 "github.com/ophum/humstack/pkg/api/system/nodenetwork/v0"
	"github.com/ophum/humstack/pkg/api/system/virtualmachine"
	vmv0 "github.com/ophum/humstack/pkg/api/system/virtualmachine/v0"
	"github.com/ophum/humstack/pkg/api/system/virtualrouter"
	vrv0 "github.com/ophum/humstack/pkg/api/system/virtualrouter/v0"
	"github.com/ophum/humstack/pkg/api/watch"
	watchv0 "github.com/ophum/humstack/pkg/api/watch/v0"
	"github.com/rakyll/statik/fs"

	//store "github.com/ophum/humstack/pkg/store/memory"
	store "github.com/ophum/humstack/pkg/store/leveldb"

	"github.com/gin-contrib/cors"
)

var (
	listenAddress string
	listenPort    int64
)

func init() {
	flag.StringVar(&listenAddress, "listen-address", "localhost", "listen address")
	flag.Int64Var(&listenPort, "listen-port", 8080, "listen port")
	flag.Parse()
}

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	notifier := make(chan string, 100)
	//s := store.NewMemoryStore()
	s, err := store.NewLevelDBStore("./database", notifier)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// bloadcasting
	notifiers := map[string](chan string){}
	go func() {
		for n := range notifier {
			for _, nn := range notifiers {
				nn <- n
			}
		}
	}()

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	r.StaticFS("/static", statikFS)

	grh := grv0.NewGroupHandler(s)
	nsh := nsv0.NewNamespaceHandler(s)
	nnwh := nodenetv0.NewNodeNetworkHandler(s)
	nwh := netv0.NewNetworkHandler(s)
	bsh := bsv0.NewBlockStorageHandler(s)
	vmh := vmv0.NewVirtualMachineHandler(s)
	vrh := vrv0.NewVirtualRouterHandler(s)
	eippoolh := eippoolv0.NewExternalIPPoolHandler(s)
	eiph := eipv0.NewExternalIPHandler(s)
	imh := imv0.NewImageHandler(s)
	ieh := iev0.NewImageEntityHandler(s)
	nodeh := nodev0.NewNodeHandler(s)
	watchh := watchv0.NewWatchHandler(notifiers)

	v0 := r.Group("/api/v0")
	{
		gri := group.NewGroupHandler(v0, grh)
		nsi := namespace.NewNamespaceHandler(v0, nsh)
		nwi := network.NewNetworkHandler(v0, nwh)
		nnwi := nodenetwork.NewNodeNetworkHandler(v0, nnwh)
		bsi := blockstorage.NewBlockStorageHandler(v0, bsh)
		vmi := virtualmachine.NewVirtualMachineHandler(v0, vmh)
		vri := virtualrouter.NewVirtualRouterHandler(v0, vrh)
		eippooli := externalippool.NewExternalIPPoolHandler(v0, eippoolh)
		eipi := externalip.NewExternalIPHandler(v0, eiph)
		imi := image.NewImageHandler(v0, imh)
		iei := imageentity.NewImageEntityHandler(v0, ieh)
		nodei := node.NewNodeHandler(v0, nodeh)
		watchi := watch.NewWatchHandler(v0, watchh)

		gri.RegisterHandlers()
		nsi.RegisterHandlers()
		nwi.RegisterHandlers()
		nnwi.RegisterHandlers()
		bsi.RegisterHandlers()
		vmi.RegisterHandlers()
		vri.RegisterHandlers()
		eippooli.RegisterHandlers()
		eipi.RegisterHandlers()
		imi.RegisterHandlers()
		iei.RegisterHandlers()
		nodei.RegisterHandlers()
		watchi.RegisterHandlers()
	}

	if err := r.Run(fmt.Sprintf("%s:%d", listenAddress, listenPort)); err != nil {
		log.Fatal(err)
	}

}
