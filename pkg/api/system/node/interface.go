package node

import (
	"github.com/gin-gonic/gin"
)

type NodeHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "nodes"
)

type NodeHandler struct {
	router *gin.RouterGroup
	nhi    NodeHandlerInterface
}

func NewNodeHandler(router *gin.RouterGroup, nhi NodeHandlerInterface) *NodeHandler {
	return &NodeHandler{
		router: router,
		nhi:    nhi,
	}
}

func (h *NodeHandler) RegisterHandlers() {
	node := h.router.Group(basePath)
	{
		node.GET("", h.nhi.FindAll)
		node.GET("/:node_name", h.nhi.Find)
		node.POST("", h.nhi.Create)
		node.PUT("/:node_name", h.nhi.Update)
		node.DELETE("/:node_name", h.nhi.Delete)
	}
}
