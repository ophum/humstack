package namespace

import (
	"github.com/gin-gonic/gin"
)

type NamespaceHandlerInterface interface {
	Find(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type NamespaceHandler struct {
	router   *gin.RouterGroup
	basePath string
	nhi      NamespaceHandlerInterface
}

func NewNamespaceHandler(router *gin.RouterGroup, basePath string, nhi NamespaceHandlerInterface) *NamespaceHandler {
	return &NamespaceHandler{
		router:   router,
		basePath: basePath,
		nhi:      nhi,
	}
}

func (h *NamespaceHandler) RegisterHandlers() {
	ns := h.router.Group(h.basePath)
	{
		ns.GET("", h.nhi.Find)
		ns.GET("/:namespace_name", h.nhi.FindById)
		ns.POST("", h.nhi.Create)
		ns.PUT("/:namespace_name", h.nhi.Update)
		ns.DELETE("/:namespace_name", h.nhi.Delete)
	}
}
