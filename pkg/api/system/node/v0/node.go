package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/node"
	"github.com/ophum/humstack/pkg/store"
)

type NodeHandler struct {
	node.NodeHandlerInterface

	store store.Store
}

func NewNodeHandler(store store.Store) *NodeHandler {
	return &NodeHandler{
		store: store,
	}
}

func (h *NodeHandler) FindAll(ctx *gin.Context) {
	nsName := getNSName(ctx)

	list := h.store.List(getKey(nsName, ""))
	nodeList := []system.Node{}
	for _, o := range list {
		nodeList = append(nodeList, o.(system.Node))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodes": nodeList,
	})

}

func (h *NodeHandler) Find(ctx *gin.Context) {
	nsName := getNSName(ctx)
	nodeName := getNodeName(ctx)

	obj := h.store.Get(getKey(nsName, nodeName))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Node `%s` is not found.", nodeName), nil)
		return
	}

	node := obj.(system.Node)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"node": node,
	})
}

func (h *NodeHandler) Create(ctx *gin.Context) {
	nsName := getNSName(ctx)

	var request system.Node
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
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Node `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"node": request,
	})
}

func (h *NodeHandler) Update(ctx *gin.Context) {
	nsName := getNSName(ctx)
	nodeName := getNodeName(ctx)

	var request system.Node
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if nodeName != request.Name {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Node Name."), nil)
		return
	}
	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(nsName, request.Name)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Node `%s` is not found.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"node": request,
	})
}

func (h *NodeHandler) Delete(ctx *gin.Context) {
	nsName := getNSName(ctx)
	nodeName := getNodeName(ctx)

	key := getKey(nsName, nodeName)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"node": nil,
	})
}

func getNSName(ctx *gin.Context) string {
	return ctx.Param("namespace_name")
}

func getNodeName(ctx *gin.Context) string {
	return ctx.Param("node_name")
}

func getKey(nsName, name string) string {
	return filepath.Join("node", nsName, name)
}
