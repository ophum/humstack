package v0

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
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

	var bs system.BlockStorage
	err := h.store.Get(getKey(groupID, nsID, bsID), &bs)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("BlockStorage `%s` is not found.", bsID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"blockstorage": bs,
	})
}

func (h *BlockStorageHandler) Create(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	var ns core.Namespace
	err := h.store.Get(filepath.Join("namespace", groupID, nsID), &ns)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: namespace is not found."), nil)
		return
	}

	var request system.BlockStorage
	err = ctx.Bind(&request)
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
	var bs system.BlockStorage
	err = h.store.Get(key, &bs)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: BlockStorage `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeBlockStorageV0
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

func (h *BlockStorageHandler) UpdateStatus(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	var request system.BlockStorage
	if err := ctx.Bind(&request); err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	key := getKey(groupID, nsID, request.ID)

	h.store.Lock(key)
	defer h.store.Unlock(key)

	var bs system.BlockStorage
	if err := h.store.Get(key, &bs); err != nil {
		meta.ResponseJSON(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	bs.Status = request.Status
	h.store.Put(key, bs)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"blockstorage": bs,
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

func (h *BlockStorageHandler) ProxyDownloadAPI(ctx *gin.Context) {
	groupID, nsID, bsID := getIDs(ctx)

	key := getKey(groupID, nsID, bsID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	var bs system.BlockStorage
	if err := h.store.Get(key, &bs); err != nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("blockstorage not found"), gin.H{})
		return
	}

	target, ok := bs.Annotations["bs-download-host"]
	if !ok {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("proxy target not found"), gin.H{})
		return
	}

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.Host = target
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
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
