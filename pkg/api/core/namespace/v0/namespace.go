package v0

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/core/namespace"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/store"
)

type NamespaceHandler struct {
	namespace.NamespaceHandlerInterface

	store store.Store
}

func NewNamespaceHandler(store store.Store) *NamespaceHandler {
	return &NamespaceHandler{
		store: store,
	}
}

func (h *NamespaceHandler) FindAll(ctx *gin.Context) {
	groupID := getGroupID(ctx)

	nsList := []*core.Namespace{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			ns := &core.Namespace{}
			nsList = append(nsList, ns)
			m = append(m, ns)
		}
		return m
	}
	h.store.List(getKey(groupID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"namespaces": nsList,
	})
}

func (h *NamespaceHandler) Find(ctx *gin.Context) {
	groupID := getGroupID(ctx)
	nsID := getNSID(ctx)
	obj := h.store.Get(getKey(groupID, nsID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Namespace `%s` is not found.", nsID), nil)
		return
	}

	ns := obj.(core.Namespace)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"namespace": ns,
	})
}

func (h *NamespaceHandler) Create(ctx *gin.Context) {
	var request core.Namespace

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

	key := getKey(request.Group, request.ID)
	obj := h.store.Get(key)
	if obj != nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: namespace `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"namespace": request,
	})
}

func (h *NamespaceHandler) Update(ctx *gin.Context) {
	groupID := getGroupID(ctx)
	nsID := getNSID(ctx)

	var request core.Namespace

	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID != nsID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: can't change id."), nil)
		return
	}

	if request.Group != groupID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: can't change group"), nil)
	}

	key := getKey(request.Group, request.ID)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Error: namespace `%s` is not found.", request.ID), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"namespace": request,
	})
}

func (h *NamespaceHandler) Delete(ctx *gin.Context) {
	groupID := getGroupID(ctx)
	nsID := getNSID(ctx)

	key := getKey(groupID, nsID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
}

func getGroupID(ctx *gin.Context) string {
	return ctx.Param("group_id")
}

func getNSID(ctx *gin.Context) string {
	return ctx.Param("namespace_id")
}

func getKey(groupID, nsID string) string {
	return filepath.Join("namespace", groupID, nsID)
}
