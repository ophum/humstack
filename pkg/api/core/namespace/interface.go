package namespace

import (
	"github.com/gin-gonic/gin"
)

type NamespaceHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "groups/:group_id/namespaces"
)

type NamespaceHandler struct {
	router *gin.RouterGroup
	nhi    NamespaceHandlerInterface
}

func NewNamespaceHandler(router *gin.RouterGroup, nhi NamespaceHandlerInterface) *NamespaceHandler {
	return &NamespaceHandler{
		router: router,
		nhi:    nhi,
	}
}

func (h *NamespaceHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.nhi.FindAll)
		ns.GET("/:namespace_id", h.nhi.Find)
		ns.POST("", h.nhi.Create)
		ns.PUT("/:namespace_id", h.nhi.Update)
		ns.DELETE("/:namespace_id", h.nhi.Delete)
	}
}
