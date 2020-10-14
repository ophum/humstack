package v0

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/core/externalip"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/store"
)

type ExternalIPHandler struct {
	externalip.ExternalIPHandlerInterface

	store store.Store
}

func NewExternalIPHandler(store store.Store) *ExternalIPHandler {
	return &ExternalIPHandler{
		store: store,
	}
}

func (h *ExternalIPHandler) FindAll(ctx *gin.Context) {

	eipList := []*core.ExternalIP{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			eip := &core.ExternalIP{}
			eipList = append(eipList, eip)
			m = append(m, eip)
		}
		return m
	}

	h.store.List(getKey("")+"/", f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalips": eipList,
	})
}

func (h *ExternalIPHandler) Find(ctx *gin.Context) {
	eipID := getExternalIPID(ctx)

	var eip core.ExternalIP
	err := h.store.Get(getKey(eipID), &eip)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("ExternalIP `%s` is not found.", eipID), nil)
		return
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalip": eip,
	})
}

func (h *ExternalIPHandler) Create(ctx *gin.Context) {
	var request core.ExternalIP

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
	var eip core.ExternalIP
	err = h.store.Get(key, &eip)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: externalip `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeExternalIPV0
	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"externalip": request,
	})
}

func (h *ExternalIPHandler) Update(ctx *gin.Context) {
	eipID := getExternalIPID((ctx))
	var request core.ExternalIP

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID != eipID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: can't change id."), nil)
		return
	}

	key := getKey(request.ID)
	var eip core.ExternalIP
	err = h.store.Get(key, &eip)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: externalip `%s` is not found.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"externalip": request,
	})
}

func (h *ExternalIPHandler) Delete(ctx *gin.Context) {
	eipID := getExternalIPID(ctx)

	key := getKey(eipID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
}

func getExternalIPID(ctx *gin.Context) string {
	return ctx.Param("external_ip_id")
}

func getKey(id string) string {
	return filepath.Join("externalip", id)
}
