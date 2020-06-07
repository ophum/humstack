package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core/namespace"
	v0 "github.com/ophum/humstack/pkg/api/core/namespace/v0"
	store "github.com/ophum/humstack/pkg/store/memory"
)

func main() {
	r := gin.Default()

	nsStore := store.NewMemoryStore()
	nsh := v0.NewNamespaceHandler(nsStore)

	v0 := r.Group("/api/v0")
	{
		nsi := namespace.NewNamespaceHandler(v0, "namespaces", nsh)
		nsi.RegisterHandlers()
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
