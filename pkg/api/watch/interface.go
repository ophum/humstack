package watch

import (
	"github.com/gin-gonic/gin"
)

type WatchHandlerInterface interface {
	Watch(ctx *gin.Context)
}

type WatchHandler struct {
	router *gin.RouterGroup
	whi    WatchHandlerInterface
}

const (
	basePath = "watches"
)

func NewWatchHandler(router *gin.RouterGroup, whi WatchHandlerInterface) *WatchHandler {
	return &WatchHandler{
		router: router,
		whi:    whi,
	}
}

func (h *WatchHandler) RegisterHandlers() {
	ns := h.router.Group(basePath)
	{
		ns.GET("", h.whi.Watch)
	}
}
