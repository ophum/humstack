package v0

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/nodenetwork"
	"github.com/ophum/humstack/pkg/store"
)

type NodeNetworkHandler struct {
	nodenetwork.NodeNetworkHandlerInterface

	store store.Store
}

func NewNodeNetworkHandler(store store.Store) *NodeNetworkHandler {
	return &NodeNetworkHandler{
		store: store,
	}
}

func (h *NodeNetworkHandler) FindAll(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	netList := []*system.NodeNetwork{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			net := &system.NodeNetwork{}
			netList = append(netList, net)
			m = append(m, net)
		}
		return m
	}

	h.store.List(getKey(groupID, nsID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodenetworks": netList,
	})
}

func (h *NodeNetworkHandler) Find(ctx *gin.Context) {
	groupID, nsID, netID := getIDs(ctx)

	var net system.NodeNetwork
	err := h.store.Get(getKey(groupID, nsID, netID), &net)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("NodeNetwork `%s` is not found.", nsID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodenetwork": net,
	})
}

func (h *NodeNetworkHandler) Create(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	var request system.NodeNetwork

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
	var net system.NodeNetwork
	err = h.store.Get(key, &net)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: NodeNetwork `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeNodeNetworkV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"nodenetwork": request,
	})
}

func (h *NodeNetworkHandler) Update(ctx *gin.Context) {
	groupID, nsID, netID := getIDs(ctx)

	var request system.NodeNetwork
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if netID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change NodeNetwork Name."), nil)
		return
	}

	key := getKey(groupID, nsID, netID)
	var net system.NodeNetwork
	err = h.store.Get(key, &net)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: NodeNetwork `%s` is not found in Namespace `%s`.", netID, nsID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodenetwork": request,
	})
}

func (h *NodeNetworkHandler) Delete(ctx *gin.Context) {
	groupID, nsID, netID := getIDs(ctx)

	key := getKey(groupID, nsID, netID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodenetwork": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, nsID, netID string) {
	groupID = ctx.Param("group_id")
	nsID = ctx.Param("namespace_id")
	netID = ctx.Param("node_network_id")
	return groupID, nsID, netID
}

func getKey(groupID, nsID, id string) string {
	return filepath.Join("nodenetwork", groupID, nsID, id)
}
