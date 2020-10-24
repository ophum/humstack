package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/imageentity"
	"github.com/ophum/humstack/pkg/store"
)

type ImageEntityHandler struct {
	imageentity.ImageEntityHandlerInterface

	store store.Store
}

func NewImageEntityHandler(store store.Store) *ImageEntityHandler {
	return &ImageEntityHandler{
		store: store,
	}
}

func (h *ImageEntityHandler) FindAll(ctx *gin.Context) {
	groupID, _ := getIDs(ctx)

	imList := []*system.ImageEntity{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			im := &system.ImageEntity{}
			imList = append(imList, im)
			m = append(m, im)
		}
		return m
	}

	h.store.List(getKey(groupID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"imageentities": imList,
	})

}

func (h *ImageEntityHandler) Find(ctx *gin.Context) {
	groupID, imID := getIDs(ctx)

	var im system.ImageEntity
	err := h.store.Get(getKey(groupID, imID), &im)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("ImageEntity `%s` is not found.", imID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"imageentity": im,
	})
}

func (h *ImageEntityHandler) Create(ctx *gin.Context) {
	groupID, _ := getIDs(ctx)

	var request system.ImageEntity
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: id is empty."), nil)
		return
	}

	key := getKey(groupID, request.ID)
	var im system.ImageEntity
	err = h.store.Get(key, &im)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: ImageEntity `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeImageEntityV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"imageentity": request,
	})
}

func (h *ImageEntityHandler) Update(ctx *gin.Context) {
	groupID, imID := getIDs(ctx)

	var request system.ImageEntity
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if imID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change ImageEntity Name."), nil)
		return
	}
	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(groupID, request.ID)

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"imageentity": request,
	})
}

func (h *ImageEntityHandler) Delete(ctx *gin.Context) {
	groupID, imID := getIDs(ctx)

	key := getKey(groupID, imID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"imageentity": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, imID string) {
	groupID = ctx.Param("group_id")
	imID = ctx.Param("image_entity_id")
	return groupID, imID
}

func getKey(groupID, id string) string {
	return filepath.Join("imageentities", groupID, id)
}
