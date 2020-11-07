package v0

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/core/network"
	"github.com/ophum/humstack/pkg/api/meta"
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
	groupID, nsID, _ := getIDs(ctx)

	netList := []*core.Network{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			net := &core.Network{}
			netList = append(netList, net)
			m = append(m, net)
		}
		return m
	}

	h.store.List(getKey(groupID, nsID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"networks": netList,
	})
}

func (h *NetworkHandler) Find(ctx *gin.Context) {
	groupID, nsID, netID := getIDs(ctx)

	var net core.Network
	err := h.store.Get(getKey(groupID, nsID, netID), &net)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Network `%s` is not found.", nsID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": net,
	})
}

func (h *NetworkHandler) Create(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	var request core.Network

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: id is empty."), nil)
		return
	}

	var ns core.Namespace
	err = h.store.Get(filepath.Join("namespace", groupID, nsID), &ns)
	if err != nil && err.Error() == "Not Found" {
		log.Println("error")
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: namespace is not found."), nil)
		return
	}

	key := getKey(groupID, nsID, request.ID)
	var net core.Network
	err = h.store.Get(key, &net)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Network `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeNetworkV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"network": request,
	})
}

func (h *NetworkHandler) Update(ctx *gin.Context) {
	groupID, nsID, netID := getIDs(ctx)

	var request core.Network
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if netID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Network Name."), nil)
		return
	}

	key := getKey(groupID, nsID, netID)
	var net core.Network
	err = h.store.Get(key, &net)
	if err != nil && err.Error() == "Not Found" {
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
	groupID, nsID, netID := getIDs(ctx)

	key := getKey(groupID, nsID, netID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"network": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, nsID, netID string) {
	groupID = ctx.Param("group_id")
	nsID = ctx.Param("namespace_id")
	netID = ctx.Param("network_id")
	return groupID, nsID, netID
}

func getKey(groupID, nsID, id string) string {
	return filepath.Join("network", groupID, nsID, id)
}
