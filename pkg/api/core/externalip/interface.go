package externalip

import (
	"github.com/gin-gonic/gin"
)

type ExternalIPHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "externalips"
)

type ExternalIPHandler struct {
	router *gin.RouterGroup
	eipi   ExternalIPHandlerInterface
}

func NewExternalIPHandler(router *gin.RouterGroup, eipi ExternalIPHandlerInterface) *ExternalIPHandler {
	return &ExternalIPHandler{
		router: router,
		eipi:   eipi,
	}
}

func (h *ExternalIPHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.eipi.FindAll)
		ns.GET("/:external_ip_id", h.eipi.Find)
		ns.POST("", h.eipi.Create)
		ns.PUT("/:external_ip_id", h.eipi.Update)
		ns.DELETE("/:external_ip_id", h.eipi.Delete)
	}
}
