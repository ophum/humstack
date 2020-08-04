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
	"github.com/ophum/humstack/pkg/api/system/virtualrouter"
	"github.com/ophum/humstack/pkg/store"
)

type VirtualRouterHandler struct {
	virtualrouter.VirtualRouterHandlerInterface

	store store.Store
}

func NewVirtualRouterHandler(store store.Store) *VirtualRouterHandler {
	return &VirtualRouterHandler{
		store: store,
	}
}

func (h *VirtualRouterHandler) FindAll(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	vrList := []*system.VirtualRouter{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			vr := &system.VirtualRouter{}
			vrList = append(vrList, vr)
			m = append(m, vr)
		}
		return m
	}

	h.store.List(getKey(groupID, nsID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualrouters": vrList,
	})
}

func (h *VirtualRouterHandler) Find(ctx *gin.Context) {
	groupID, nsID, vrID := getIDs(ctx)

	var vr system.VirtualRouter
	err := h.store.Get(getKey(groupID, nsID, vrID), &vr)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("VirtualRouter `%s` is not found.", nsID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualrouter": vr,
	})
}

func (h *VirtualRouterHandler) Create(ctx *gin.Context) {
	groupID, nsID, _ := getIDs(ctx)

	var request system.VirtualRouter

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
	var vr system.VirtualRouter
	err = h.store.Get(key, &vr)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: VirtualRouter `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"virtualrouter": request,
	})
}

func (h *VirtualRouterHandler) Update(ctx *gin.Context) {
	groupID, nsID, vrID := getIDs(ctx)

	var request system.VirtualRouter
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if vrID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change VirtualRouter Name."), nil)
		return
	}

	key := getKey(groupID, nsID, vrID)
	var vr system.VirtualRouter
	err = h.store.Get(key, &vr)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: VirtualRouter `%s` is not found in Namespace `%s`.", vrID, nsID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualrouter": request,
	})
}

func (h *VirtualRouterHandler) Delete(ctx *gin.Context) {
	groupID, nsID, vrID := getIDs(ctx)

	key := getKey(groupID, nsID, vrID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualrouter": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, nsID, vrID string) {
	groupID = ctx.Param("group_id")
	nsID = ctx.Param("namespace_id")
	vrID = ctx.Param("virtualrouter_id")
	return groupID, nsID, vrID
}

func getKey(groupID, nsID, id string) string {
	return filepath.Join("virtualrouter", groupID, nsID, id)
}
