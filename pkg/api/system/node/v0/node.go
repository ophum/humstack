package v0

import (
	"fmt"
	"log"
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

	nodeList := []*system.Node{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			node := &system.Node{}
			nodeList = append(nodeList, node)
			m = append(m, node)
		}
		return m
	}

	h.store.List(getKey("")+"/", f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodes": nodeList,
	})

}

func (h *NodeHandler) Find(ctx *gin.Context) {
	nodeID := getNodeID(ctx)

	var node system.Node
	err := h.store.Get(getKey(nodeID), &node)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Node `%s` is not found.", nodeID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"node": node,
	})
}

func (h *NodeHandler) Create(ctx *gin.Context) {
	var request system.Node
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	err = h.validate(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	key := getKey(request.ID)
	var node system.Node
	err = h.store.Get(key, &node)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: Node `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeNodeV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"node": request,
	})
}

func (h *NodeHandler) Update(ctx *gin.Context) {
	nodeID := getNodeID(ctx)

	var request system.Node
	err := ctx.Bind(&request)
	if err != nil {
		log.Println(err)
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if nodeID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Node ID."), nil)
		return
	}

	err = h.validate(&request)
	if err != nil {
		if err.Error() != "Error: id is duplicated." {
			log.Println(err)
			meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
			return
		}
	}

	key := getKey(nodeID)
	var node system.Node
	err = h.store.Get(key, &node)
	if err != nil && err.Error() == "Not Found" {
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
	nodeID := getNodeID(ctx)

	key := getKey(nodeID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"node": nil,
	})
}

func (h *NodeHandler) isIDDuplicate(node *system.Node) bool {
	list := []*system.Node{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			node := &system.Node{}
			list = append(list, node)
			m = append(m, node)
		}
		return m
	}
	h.store.List(getKey(""), f)
	for _, n := range list {
		if n.ID == node.ID {
			return true
		}
	}
	return false
}

func (h *NodeHandler) validate(node *system.Node) error {
	if node.ID == "" {
		return fmt.Errorf("Error: id is empty.")
	}

	if h.isIDDuplicate(node) {
		return fmt.Errorf("Error: id is duplicated.")
	}
	return nil
}

func getNodeID(ctx *gin.Context) string {
	return ctx.Param("node_id")
}

func getKey(name string) string {
	return filepath.Join("node", name)
}
