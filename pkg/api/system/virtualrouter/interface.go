package virtualrouter

import (
	"github.com/gin-gonic/gin"
)

type VirtualRouterHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type VirtualRouterHandler struct {
	router *gin.RouterGroup
	vrhi   VirtualRouterHandlerInterface
}

const (
	basePath = "groups/:group_id/namespaces/:namespace_id/virtualrouters"
)

func NewVirtualRouterHandler(router *gin.RouterGroup, vrhi VirtualRouterHandlerInterface) *VirtualRouterHandler {
	return &VirtualRouterHandler{
		router: router,
		vrhi:   vrhi,
	}
}

func (h *VirtualRouterHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.vrhi.FindAll)
		ns.GET("/:virtualrouter_id", h.vrhi.Find)
		ns.POST("", h.vrhi.Create)
		ns.PUT("/:virtualrouter_id", h.vrhi.Update)
		ns.DELETE("/:virtualrouter_id", h.vrhi.Delete)
	}
}
