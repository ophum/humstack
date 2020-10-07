package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core/externalip"
	eipv0 "github.com/ophum/humstack/pkg/api/core/externalip/v0"
	"github.com/ophum/humstack/pkg/api/core/externalippool"
	eippoolv0 "github.com/ophum/humstack/pkg/api/core/externalippool/v0"
	"github.com/ophum/humstack/pkg/api/core/group"
	grv0 "github.com/ophum/humstack/pkg/api/core/group/v0"
	"github.com/ophum/humstack/pkg/api/core/namespace"
	nsv0 "github.com/ophum/humstack/pkg/api/core/namespace/v0"
	"github.com/ophum/humstack/pkg/api/system/blockstorage"
	bsv0 "github.com/ophum/humstack/pkg/api/system/blockstorage/v0"
	"github.com/ophum/humstack/pkg/api/system/network"
	netv0 "github.com/ophum/humstack/pkg/api/system/network/v0"
	"github.com/ophum/humstack/pkg/api/system/node"
	nodev0 "github.com/ophum/humstack/pkg/api/system/node/v0"
	"github.com/ophum/humstack/pkg/api/system/virtualmachine"
	vmv0 "github.com/ophum/humstack/pkg/api/system/virtualmachine/v0"
	"github.com/ophum/humstack/pkg/api/system/virtualrouter"
	vrv0 "github.com/ophum/humstack/pkg/api/system/virtualrouter/v0"

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

	//s := store.NewMemoryStore()
	s, err := store.NewLevelDBStore("./database")
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	grh := grv0.NewGroupHandler(s)
	nsh := nsv0.NewNamespaceHandler(s)
	nwh := netv0.NewNetworkHandler(s)
	bsh := bsv0.NewBlockStorageHandler(s)
	vmh := vmv0.NewVirtualMachineHandler(s)
	vrh := vrv0.NewVirtualRouterHandler(s)
	eippoolh := eippoolv0.NewExternalIPPoolHandler(s)
	eiph := eipv0.NewExternalIPHandler(s)
	nodeh := nodev0.NewNodeHandler(s)

	v0 := r.Group("/api/v0")
	{
		gri := group.NewGroupHandler(v0, grh)
		nsi := namespace.NewNamespaceHandler(v0, nsh)
		nwi := network.NewNetworkHandler(v0, nwh)
		bsi := blockstorage.NewBlockStorageHandler(v0, bsh)
		vmi := virtualmachine.NewVirtualMachineHandler(v0, vmh)
		vri := virtualrouter.NewVirtualRouterHandler(v0, vrh)
		eippooli := externalippool.NewExternalIPPoolHandler(v0, eippoolh)
		eipi := externalip.NewExternalIPHandler(v0, eiph)
		nodei := node.NewNodeHandler(v0, nodeh)

		gri.RegisterHandlers()
		nsi.RegisterHandlers()
		nwi.RegisterHandlers()
		bsi.RegisterHandlers()
		vmi.RegisterHandlers()
		vri.RegisterHandlers()
		eippooli.RegisterHandlers()
		eipi.RegisterHandlers()
		nodei.RegisterHandlers()
	}

	if err := r.Run(fmt.Sprintf("%s:%d", listenAddress, listenPort)); err != nil {
		log.Fatal(err)
	}

}
