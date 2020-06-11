package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	nsID := ctx.Param("namespace_id")

	list := h.store.List(getKey(nsID, ""))
	nsList := []system.Network{}
	for _, o := range list {
		nsList = append(nsList, o.(system.Network))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"networks": nsList,
	})
}

func (h *NetworkHandler) Find(ctx *gin.Context) {
	nsID := ctx.Param("namespace_id")
	netID := ctx.Param("network_id")

	obj := h.store.Get(getKey(nsID, netID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Network `%s` is not found.", nsID), nil)
		return
	}

	ns := obj.(system.Network)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": ns,
	})
}

func (h *NetworkHandler) Create(ctx *gin.Context) {
	nsID := ctx.Param("namespace_id")

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

	obj := h.store.Get(filepath.Join("namespace", nsID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: namespace is not found."), nil)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	request.ID = id.String()
	key := getKey(nsID, request.ID)
	obj = h.store.Get(key)
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
	nsID := ctx.Param("namespace_id")
	netID := ctx.Param("network_id")

	var request system.Network
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if netID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Network Name."), nil)
		return
	}

	key := getKey(nsID, netID)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: Network `%s` is not found in Namespace `%s`.", netID, nsID), nil)
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
	nsID := ctx.Param("namespace_id")
	netID := ctx.Param("network_id")

	key := getKey(nsID, netID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": nil,
	})
}

func getKey(nsID, id string) string {
	return filepath.Join("network", nsID, id)
}
