package group

import (
	"github.com/gin-gonic/gin"
)

type GroupHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

const (
	basePath = "groups"
)

type GroupHandler struct {
	router *gin.RouterGroup
	nhi    GroupHandlerInterface
}

func NewGroupHandler(router *gin.RouterGroup, nhi GroupHandlerInterface) *GroupHandler {
	return &GroupHandler{
		router: router,
		nhi:    nhi,
	}
}

func (h *GroupHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.nhi.FindAll)
		ns.GET("/:group_id", h.nhi.Find)
		ns.POST("", h.nhi.Create)
		ns.PUT("/:group_id", h.nhi.Update)
		ns.DELETE("/:group_id", h.nhi.Delete)
	}
}
