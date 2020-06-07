package network

import (
	"github.com/gin-gonic/gin"
)

type NetworkHandlerInterface interface {
	Find(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type NetworkHandler struct {
	router   *gin.RouterGroup
	basePath string
	nhi      NetworkHandlerInterface
}

func NewNetworkHandler(router *gin.RouterGroup, basePath string, nhi NetworkHandlerInterface) *NetworkHandler {
	return &NetworkHandler{
		router:   router,
		basePath: basePath,
		nhi:      nhi,
	}
}

func (h *NetworkHandler) RegisterHandlers() {
	ns := h.router.Group(h.basePath)
	{
		ns.GET("", h.nhi.Find)
		ns.GET("/:network_name", h.nhi.FindById)
		ns.POST("", h.nhi.Create)
		ns.PUT("/:network_name", h.nhi.Update)
		ns.DELETE("/:network_name", h.nhi.Delete)
	}
}
