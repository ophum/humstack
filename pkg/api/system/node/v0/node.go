package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/node"
	"github.com/ophum/humstack/pkg/store"
	"gopkg.in/yaml.v2"
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

	list := h.store.List(getKey(""))
	nodeList := []system.Node{}
	for _, o := range list {
		nodeList = append(nodeList, o.(system.Node))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"nodes": nodeList,
	})

}

func (h *NodeHandler) Find(ctx *gin.Context) {
	nodeID := getNodeID(ctx)

	obj := h.store.Get(getKey(nodeID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Node `%s` is not found.", nodeID), nil)
		return
	}

	node := obj.(system.Node)
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

	buf, _ := yaml.Marshal(request)

	fmt.Println(string(buf))
	err = h.validate(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusInternalServerError, fmt.Errorf("Error: failed to generate id."), nil)
		return
	}

	request.ID = id.String()

	key := getKey(request.ID)
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
	nodeID := getNodeID(ctx)

	var request system.Node
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if nodeID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change Node Name."), nil)
		return
	}

	err = h.validate(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	key := getKey(nodeID)
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
	nodeID := getNodeID(ctx)

	key := getKey(nodeID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"node": nil,
	})
}

func (h *NodeHandler) isNameDuplicate(name string) bool {
	list := h.store.List(getKey(""))
	for _, o := range list {
		if o.(system.Node).Name == name {
			return true
		}
	}
	return false
}

func (h *NodeHandler) validate(node *system.Node) error {
	if node.Name == "" {
		return fmt.Errorf("Error: name is empty.")
	}

	if h.isNameDuplicate(node.Name) {
		return fmt.Errorf("Error: name is empty.")
	}
	return nil
}

func getNodeID(ctx *gin.Context) string {
	return ctx.Param("node_id")
}

func getKey(id string) string {
	return filepath.Join("node", id)
}
