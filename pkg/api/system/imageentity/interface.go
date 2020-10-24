package imageentity

import "github.com/gin-gonic/gin"

type ImageEntityHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type ImageEntityHandler struct {
	router *gin.RouterGroup
	iehi   ImageEntityHandlerInterface
}

const (
	basePath = "groups/:group_id/imageentities"
)

func NewImageEntityHandler(router *gin.RouterGroup, iehi ImageEntityHandlerInterface) *ImageEntityHandler {
	return &ImageEntityHandler{
		router: router,
		iehi:   iehi,
	}
}

func (h *ImageEntityHandler) RegisterHandlers() {
	ie := h.router.Group(basePath)
	{
		ie.GET("", h.iehi.FindAll)
		ie.GET("/:image_entity_id", h.iehi.Find)
		ie.POST("", h.iehi.Create)
		ie.PUT("/:image_entity_id", h.iehi.Update)
		ie.DELETE("/:image_entity_id", h.iehi.Delete)
	}
}
