package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/network"
	"github.com/ophum/humstack/pkg/store"
)

type NetworkHandler struct {
	network.NetworkHandlerInterface

	store store.Store
}

func NewNetworkHandler(store store.Store) *NetworkHandler {
	return &NetworkHandler{
		store: store,
	}
}

func (h *NetworkHandler) FindAll(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")

	list := h.store.List(getKey(nsName, ""))
	nsList := []system.Network{}
	for _, o := range list {
		nsList = append(nsList, o.(system.Network))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"networks": nsList,
	})
}

func (h *NetworkHandler) Find(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")
	netName := ctx.Param("network_name")

	obj := h.store.Get(getKey(nsName, netName))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Network `%s` is not found.", nsName), nil)
		return
	}

	ns := obj.(system.Network)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": ns,
	})
}

func (h *NetworkHandler) Create(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")

	var request system.Network

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(nsName, request.Name)
	obj := h.store.Get(key)
	if obj != nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Network `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"network": request,
	})
}

func (h *NetworkHandler) Update(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")
	netName := ctx.Param("network_name")

	var request system.Network
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if netName != request.Name {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Network Name."), nil)
		return
	}

	key := getKey(nsName, netName)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: Network `%s` is not found in Namespace `%s`.", netName, nsName), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": request,
	})
}

func (h *NetworkHandler) Delete(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")
	netName := ctx.Param("network_name")

	key := getKey(nsName, netName)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": nil,
	})
}

func getKey(nsName, name string) string {
	return filepath.Join("network", nsName, name)
}
