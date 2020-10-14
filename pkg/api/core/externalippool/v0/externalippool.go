package v0

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/core/externalippool"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/store"
)

type ExternalIPPoolHandler struct {
	externalippool.ExternalIPPoolHandlerInterface

	store store.Store
}

func NewExternalIPPoolHandler(store store.Store) *ExternalIPPoolHandler {
	return &ExternalIPPoolHandler{
		store: store,
	}
}

func (h *ExternalIPPoolHandler) FindAll(ctx *gin.Context) {

	eippoolList := []*core.ExternalIPPool{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			eippool := &core.ExternalIPPool{}
			eippoolList = append(eippoolList, eippool)
			m = append(m, eippool)
		}
		return m
	}

	h.store.List(getKey("")+"/", f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalippools": eippoolList,
	})
}

func (h *ExternalIPPoolHandler) Find(ctx *gin.Context) {
	eippoolID := getExternalIPPoolID(ctx)

	var eippool core.ExternalIPPool
	err := h.store.Get(getKey(eippoolID), &eippool)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("ExternalIPPool `%s` is not found.", eippoolID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalippool": eippool,
	})
}

func (h *ExternalIPPoolHandler) Create(ctx *gin.Context) {
	var request core.ExternalIPPool

	err := ctx.Bind(&request)
	if err != nil {
		log.Println(err)
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		log.Println("id is empty")
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: ID is empty."), nil)
		return
	}

	key := getKey(request.ID)
	var eippool core.ExternalIPPool
	err = h.store.Get(key, &eippool)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: externalippool `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeExternalIPPoolV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"externalippool": request,
	})
}

func (h *ExternalIPPoolHandler) Update(ctx *gin.Context) {
	eippoolID := getExternalIPPoolID((ctx))
	var request core.ExternalIPPool

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID != eippoolID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: can't change id."), nil)
		return
	}

	key := getKey(request.ID)
	var eippool core.ExternalIPPool
	err = h.store.Get(key, &eippool)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: externalippool `%s` is not found.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalippool": request,
	})
}

func (h *ExternalIPPoolHandler) Delete(ctx *gin.Context) {
	eippoolID := getExternalIPPoolID(ctx)

	key := getKey(eippoolID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
}

func getExternalIPPoolID(ctx *gin.Context) string {
	return ctx.Param("external_ip_pool_id")
}

func getKey(id string) string {
	return filepath.Join("externalippool", id)
}
