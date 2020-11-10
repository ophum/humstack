package v0

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/koding/websocketproxy"
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

	var vm system.VirtualMachine
	err := h.store.Get(getKey(groupID, nsID, vmID), &vm)
	if err != nil && err.Error() == "Not Found" {
		meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("VirtualMachine `%s` is not found.", vmID), nil)
		return
	}

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
	var vm system.VirtualMachine
	err = h.store.Get(key, &vm)
	if err == nil {
		meta.ResponseJSON(ctx, http.StatusConflict, fmt.Errorf("Error: VirtualMachine `%s` is already exists.", request.Name), nil)
		return
	}

	h.store.Lock(key)
	defer h.store.Unlock(key)

	request.APIType = meta.APITypeVirtualMachineV0
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
	var vm system.VirtualMachine
	err = h.store.Get(key, &vm)
	if err != nil && err.Error() == "Not Found" {
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

func (h *VirtualMachineHandler) OpenConsole(ctx *gin.Context) {
	groupID, nsID, vmID := getIDs(ctx)
	ctx.Redirect(307, fmt.Sprintf("/static/vnc.html?path=api/v0/groups/%s/namespaces/%s/virtualmachines/%s/ws", groupID, nsID, vmID))
}

func (h *VirtualMachineHandler) ConsoleWebSocketProxy(ctx *gin.Context) {
	groupID, nsID, vmID := getIDs(ctx)

	key := getKey(groupID, nsID, vmID)

	vm := system.VirtualMachine{}
	if err := h.store.Get(key, &vm); err != nil {
		if err.Error() == "Not Found" {
			meta.ResponseJSON(ctx, http.StatusNotFound, fmt.Errorf("VirtualMachine `%s` is not found.", vmID), nil)
			return
		} else {
			meta.ResponseJSON(ctx, http.StatusInternalServerError, err, gin.H{})
			log.Println(err.Error())
			return
		}
	}

	backendHost, ok := vm.Annotations["virtualmachinev0/vnc_websocket_host"]
	if !ok {
		meta.ResponseJSON(ctx, http.StatusInternalServerError, fmt.Errorf("vnc setting not found"), gin.H{})
		return
	}

	backendURL := &url.URL{
		Scheme: "ws",
		Host:   backendHost,
		Path:   "/",
	}

	ws := &websocketproxy.WebsocketProxy{
		Backend: func(*http.Request) *url.URL {
			return backendURL
		},
	}
	delete(ctx.Request.Header, "Origin")
	ws.ServeHTTP(ctx.Writer, ctx.Request)
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
