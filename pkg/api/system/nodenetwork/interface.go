package nodenetwork

import (
	"github.com/gin-gonic/gin"
)

type NodeNetworkHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type NodeNetworkHandler struct {
	router *gin.RouterGroup
	nhi    NodeNetworkHandlerInterface
}

const (
	basePath = "groups/:group_id/namespaces/:namespace_id/nodenetworks"
)

func NewNodeNetworkHandler(router *gin.RouterGroup, nhi NodeNetworkHandlerInterface) *NodeNetworkHandler {
	return &NodeNetworkHandler{
		router: router,
		nhi:    nhi,
	}
}

func (h *NodeNetworkHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.nhi.FindAll)
		ns.GET("/:node_network_id", h.nhi.Find)
		ns.POST("", h.nhi.Create)
		ns.PUT("/:node_network_id", h.nhi.Update)
		ns.DELETE("/:node_network_id", h.nhi.Delete)
	}
}
