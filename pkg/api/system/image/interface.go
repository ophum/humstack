package image

import "github.com/gin-gonic/gin"

type ImageHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	ProxyDownloadAPI(ctx *gin.Context)
}

type ImageHandler struct {
	router *gin.RouterGroup
	imhi   ImageHandlerInterface
}

const (
	basePath = "groups/:group_id/images"
)

func NewImageHandler(router *gin.RouterGroup, imhi ImageHandlerInterface) *ImageHandler {
	return &ImageHandler{
		router: router,
		imhi:   imhi,
	}
}

func (h *ImageHandler) RegisterHandlers() {
	im := h.router.Group(basePath)
	{
		im.GET("", h.imhi.FindAll)
		im.GET("/:image_id", h.imhi.Find)
		im.POST("", h.imhi.Create)
		im.PUT("/:image_id", h.imhi.Update)
		im.DELETE("/:image_id", h.imhi.Delete)
		im.GET("/:image_id/tags/:tag/download", h.imhi.ProxyDownloadAPI)
	}
}
