package externalippool

import (
	"github.com/gin-gonic/gin"
)

type ExternalIPPoolHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "externalippools"
)

type ExternalIPPoolHandler struct {
	router *gin.RouterGroup
	eipool ExternalIPPoolHandlerInterface
}

func NewExternalIPPoolHandler(router *gin.RouterGroup, eipool ExternalIPPoolHandlerInterface) *ExternalIPPoolHandler {
	return &ExternalIPPoolHandler{
		router: router,
		eipool: eipool,
	}
}

func (h *ExternalIPPoolHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.eipool.FindAll)
		ns.GET("/:external_ip_pool_id", h.eipool.Find)
		ns.POST("", h.eipool.Create)
		ns.PUT("/:external_ip_pool_id", h.eipool.Update)
		ns.DELETE("/:external_ip_pool_id", h.eipool.Delete)
	}
}
