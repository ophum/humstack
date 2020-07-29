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
	groupID, nsID, _ := getIDs(ctx)

	vmList := []*system.VirtualMachine{}
	f := func(n int) []interface{} {
		m := []interface{}{}
		for i := 0; i < n; i++ {
			vm := &system.VirtualMachine{}
			vmList = append(vmList, vm)
			m = append(m, vm)
		}
		return m
	}
	h.store.List(getKey(groupID, nsID, ""), f)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachines": vmList,
	})

}

func (h *VirtualMachineHandler) Find(ctx *gin.Context) {
	groupID, nsID, vmID := getIDs(ctx)

	obj := h.store.Get(getKey(groupID, nsID, vmID))
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
	groupID, nsID, _ := getIDs(ctx)

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

	key := getKey(groupID, nsID, request.ID)
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
	groupID, nsID, vmID := getIDs(ctx)

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

	key := getKey(groupID, nsID, request.ID)
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
	groupID, nsID, vmID := getIDs(ctx)

	key := getKey(groupID, nsID, vmID)
	h.store.Lock(key)
	defer h.store.Unlock(key)

	h.store.Delete(key)

	meta.ResponseJSON(ctx, http.StatusOK, nil, gin.H{
		"virtualmachine": nil,
	})
}

func getIDs(ctx *gin.Context) (groupID, nsID, vmID string) {
	groupID = ctx.Param("group_id")
	nsID = ctx.Param("namespace_id")
	vmID = ctx.Param("virtual_machine_id")
	return groupID, nsID, vmID
}

func getKey(groupID, nsID, vmID string) string {
	return filepath.Join("virtualmachine", groupID, nsID, vmID)
}
