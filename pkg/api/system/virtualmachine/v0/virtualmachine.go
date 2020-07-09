package v0

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
	"github.com/ophum/humstack/pkg/api/system/virtualmachine"
	"github.com/ophum/humstack/pkg/store"
)

type VirtualMachineHandler struct {
	virtualmachine.VirtualMachineHandlerInterface

	store store.Store
}

func NewVirtualMachineHandler(store store.Store) *VirtualMachineHandler {
	return &VirtualMachineHandler{
		store: store,
	}
}

func (h *VirtualMachineHandler) FindAll(ctx *gin.Context) {
	nsID := getNSID(ctx)

	list := h.store.List(getKey(nsID, ""))
	vmList := []system.VirtualMachine{}
	for _, o := range list {
		vmList = append(vmList, o.(system.VirtualMachine))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachines": vmList,
	})

}

func (h *VirtualMachineHandler) Find(ctx *gin.Context) {
	nsID := getNSID(ctx)
	vmID := getVMID(ctx)

	obj := h.store.Get(getKey(nsID, vmID))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("VirtualMachine `%s` is not found.", vmID), nil)
		return
	}

	vm := obj.(system.VirtualMachine)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachine": vm,
	})
}

func (h *VirtualMachineHandler) Create(ctx *gin.Context) {
	nsID := getNSID(ctx)

	var request system.VirtualMachine
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.ID == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: id is empty."), nil)
		return
	}

	key := getKey(nsID, request.ID)
	obj := h.store.Get(key)
	if obj != nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: VirtualMachine `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"virtualmachine": request,
	})
}

func (h *VirtualMachineHandler) Update(ctx *gin.Context) {
	nsID := getNSID(ctx)
	vmID := getVMID(ctx)

	var request system.VirtualMachine
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if vmID != request.ID {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change VirtualMachine Name."), nil)
		return
	}

	key := getKey(nsID, request.ID)
	obj := h.store.Get(key)
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: VirtualMachine `%s` is not found.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Put(key, request)

	meta.ResponseJSON(ctx, http.StatusCreated, nil, gin.H{
		"virtualmachine": request,
	})
}

func (h *VirtualMachineHandler) Delete(ctx *gin.Context) {
	nsID := getNSID(ctx)
	vmID := getVMID(ctx)

	key := getKey(nsID, vmID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachine": nil,
	})
}

func getNSID(ctx *gin.Context) string {
	return ctx.Param("namespace_id")
}

func getVMID(ctx *gin.Context) string {
	return ctx.Param("virtual_machine_id")
}

func getKey(nsID, vmID string) string {
	return filepath.Join("virtualmachine", nsID, vmID)
}
