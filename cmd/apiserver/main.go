package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core/namespace"
	nsv0 "github.com/ophum/humstack/pkg/api/core/namespace/v0"
	"github.com/ophum/humstack/pkg/api/system/blockstorage"
	bsv0 "github.com/ophum/humstack/pkg/api/system/blockstorage/v0"
	"github.com/ophum/humstack/pkg/api/system/network"
	netv0 "github.com/ophum/humstack/pkg/api/system/network/v0"
	store "github.com/ophum/humstack/pkg/store/memory"
)

func main() {
	r := gin.Default()

	s := store.NewMemoryStore()
	nsh := nsv0.NewNamespaceHandler(s)
	nwh := netv0.NewNetworkHandler(s)
	bsh := bsv0.NewBlockStorageHandler(s)

	v0 := r.Group("/api/v0")
	{
		nsi := namespace.NewNamespaceHandler(v0, nsh)
		nwi := network.NewNetworkHandler(v0, nwh)
		bsi := blockstorage.NewBlockStorageHandler(v0, bsh)
		nsi.RegisterHandlers()
		nwi.RegisterHandlers()
		bsi.RegisterHandlers()
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
