package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/blockstorage"
	"github.com/ophum/humstack/pkg/store"
)

type BlockStorageHandler struct {
	blockstorage.BlockStorageHandlerInterface

	store store.Store
}

func NewBlockStorageHandler(store store.Store) *BlockStorageHandler {
	return &BlockStorageHandler{
		store: store,
	}
}

func (h *BlockStorageHandler) FindAll(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	bsList := []*system.BlockStorage{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			bs := &system.BlockStorage{}
			bsList = append(bsList, bs)
			m = append(m, bs)
		}
		return m
	}

	h.store.List(getKey(groupID, nsID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorages": bsList,
	})

}

func (h *BlockStorageHandler) Find(ctx *gin.Context) {
	groupID, nsID, bsID := getIDs(ctx)

	obj := h.store.Get(getKey(groupID, nsID, bsID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("BlockStorage `%s` is not found.", bsID), nil)
		return
	}

	bs := obj.(system.BlockStorage)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorage": bs,
	})
}

func (h *BlockStorageHandler) Create(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	obj := h.store.Get(filepath.Join("namespace", groupID, nsID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: namespace is not found."), nil)
		return
	}

	var request system.BlockStorage
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: id is empty."), nil)
		return
	}
	if request.Spec.RequestSize == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: requestSize is empty."), nil)
		return
	}
	if request.Spec.LimitSize == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: limitSize is empty."), nil)
		return
	}

	key := getKey(groupID, nsID, request.ID)
	obj = h.store.Get(key)
	if obj != nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: BlockStorage `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"blockstorage": request,
	})
}

func (h *BlockStorageHandler) Update(ctx *gin.Context) {
	groupID, nsID, bsID := getIDs(ctx)

	var request system.BlockStorage
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if bsID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change BlockStorage Name."), nil)
		return
	}
	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}
	if request.Spec.RequestSize == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: requestSize is empty."), nil)
		return
	}
	if request.Spec.LimitSize == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: limitSize is empty."), nil)
		return
	}

	key := getKey(groupID, nsID, request.ID)

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"blockstorage": request,
	})
}

func (h *BlockStorageHandler) Delete(ctx *gin.Context) {
	groupID, nsID, bsID := getIDs(ctx)

	key := getKey(groupID, nsID, bsID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorage": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, nsID, bsID string) {
	groupID = ctx.Param("group_id")
	nsID = ctx.Param("namespace_id")
	bsID = ctx.Param("block_storage_id")
	return groupID, nsID, bsID
}

func getKey(groupID, nsID, id string) string {
	return filepath.Join("blockstorage", groupID, nsID, id)
}
