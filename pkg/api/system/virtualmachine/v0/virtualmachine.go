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
	nsName := getNSName(ctx)

	list := h.store.List(getKey(nsName, ""))
	vmList := []system.VirtualMachine{}
	for _, o := range list {
		vmList = append(vmList, o.(system.VirtualMachine))
	}

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachines": vmList,
	})

}

func (h *VirtualMachineHandler) Find(ctx *gin.Context) {
	nsName := getNSName(ctx)
	vmName := getVMName(ctx)

	obj := h.store.Get(getKey(nsName, vmName))
	if obj == nil {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("VirtualMachine `%s` is not found.", vmName), nil)
		return
	}

	vm := obj.(system.VirtualMachine)
	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachine": vm,
	})
}

func (h *VirtualMachineHandler) Create(ctx *gin.Context) {
	nsName := getNSName(ctx)

	var request system.VirtualMachine
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(nsName, request.Name)
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
	nsName := getNSName(ctx)
	vmName := getVMName(ctx)

	var request system.VirtualMachine
	err := ctx.Bind(&request)
	if err != nil {
		meta.ResponseJSON(ctx, http.StatusBadRequest, err, nil)
		return
	}

	if vmName != request.Name {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: Can't change VirtualMachine Name."), nil)
		return
	}
	if request.Name == "" {
		meta.ResponseJSON(ctx, http.StatusBadRequest, fmt.Errorf("Error: name is empty."), nil)
		return
	}

	key := getKey(nsName, request.Name)
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
	nsName := getNSName(ctx)
	vmName := getVMName(ctx)

	key := getKey(nsName, vmName)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachine": nil,
	})
}

func getNSName(ctx *gin.Context) string {
	return ctx.Param("namespace_name")
}

func getVMName(ctx *gin.Context) string {
	return ctx.Param("virtual_machine_name")
}

func getKey(nsName, name string) string {
	return filepath.Join("virtualmachine", nsName, name)
}
