package main

import (
	"log"

	"github.com/gin-gonic/gin"
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
	store "github.com/ophum/humstack/pkg/store/memory"

	"github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	s := store.NewMemoryStore()
	grh := grv0.NewGroupHandler(s)
	nsh := nsv0.NewNamespaceHandler(s)
	nwh := netv0.NewNetworkHandler(s)
	bsh := bsv0.NewBlockStorageHandler(s)
	vmh := vmv0.NewVirtualMachineHandler(s)
	nodeh := nodev0.NewNodeHandler(s)

	v0 := r.Group("/api/v0")
	{
		gri := group.NewGroupHandler(v0, grh)
		nsi := namespace.NewNamespaceHandler(v0, nsh)
		nwi := network.NewNetworkHandler(v0, nwh)
		bsi := blockstorage.NewBlockStorageHandler(v0, bsh)
		vmi := virtualmachine.NewVirtualMachineHandler(v0, vmh)
		nodei := node.NewNodeHandler(v0, nodeh)
		gri.RegisterHandlers()
		nsi.RegisterHandlers()
		nwi.RegisterHandlers()
		bsi.RegisterHandlers()
		vmi.RegisterHandlers()
		nodei.RegisterHandlers()
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
