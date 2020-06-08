package v0

import (
	"fmt"
	"net/http"

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
	list := h.store.List("namespace/")
	nsList := []core.Namespace{}
	for _, o := range list {
		nsList = append(nsList, o.(core.Namespace))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"namespaces": nsList,
	})
}

func (h *NamespaceHandler) Find(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")
	obj := h.store.Get(getKey(nsName))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("Namespace `%s` is not found.", nsName), nil)
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
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(request.Name)
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
	meta.ResponseJSON(ctx, http.StatusNotImplemented, fmt.Errorf("Not implemented"), nil)
}

func (h *NamespaceHandler) Delete(ctx *gin.Context) {
	nsName := ctx.Param("namespace_name")

	key := getKey(nsName)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)
}

func getKey(name string) string {
	return "namespace/" + name
}
