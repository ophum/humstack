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
	nsName := getNSName(ctx)

	list := h.store.List(getKey(nsName, ""))
	bsList := []system.BlockStorage{}
	for _, o := range list {
		bsList = append(bsList, o.(system.BlockStorage))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorages": bsList,
	})

}

func (h *BlockStorageHandler) Find(ctx *gin.Context) {
	nsName := getNSName(ctx)
	bsName := getBSName(ctx)

	obj := h.store.Get(getKey(nsName, bsName))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("BlockStorage `%s` is not found.", bsName), nil)
		return
	}

	bs := obj.(system.BlockStorage)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorage": bs,
	})
}

func (h *BlockStorageHandler) Create(ctx *gin.Context) {
	nsName := getNSName(ctx)

	var request system.BlockStorage
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
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

	key := getKey(nsName, request.Name)
	obj := h.store.Get(key)
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
	nsName := getNSName(ctx)
	bsName := getBSName(ctx)

	var request system.BlockStorage
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if bsName != request.Name {
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

	key := getKey(nsName, request.Name)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: BlockStorage `%s` is not found.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"blockstorage": request,
	})
}

func (h *BlockStorageHandler) Delete(ctx *gin.Context) {
	nsName := getNSName(ctx)
	bsName := getBSName(ctx)

	key := getKey(nsName, bsName)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorage": nil,
	})
}

func getNSName(ctx *gin.Context) string {
	return ctx.Param("namespace_name")
}

func getBSName(ctx *gin.Context) string {
	return ctx.Param("block_storage_name")
}

func getKey(nsName, name string) string {
	return filepath.Join("blockstorage", nsName, name)
}
