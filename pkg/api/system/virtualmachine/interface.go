package virtualmachine

import (
	"github.com/gin-gonic/gin"
)

type VirtualMachineHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "namespaces/:namespace_id/virtualmachines"
)

type VirtualMachineHandler struct {
	router *gin.RouterGroup
	vmhi   VirtualMachineHandlerInterface
}

func NewVirtualMachineHandler(router *gin.RouterGroup, vmhi VirtualMachineHandlerInterface) *VirtualMachineHandler {
	return &VirtualMachineHandler{
		router: router,
		vmhi:   vmhi,
	}
}

func (h *VirtualMachineHandler) RegisterHandlers() {
	vm := h.router.Group(basePath)
	{
		vm.GET("", h.vmhi.FindAll)
		vm.GET("/:virtual_machine_id", h.vmhi.Find)
		vm.POST("", h.vmhi.Create)
		vm.PUT("/:virtual_machine_id", h.vmhi.Update)
		vm.DELETE("/:virtual_machine_id", h.vmhi.Delete)
	}
}
