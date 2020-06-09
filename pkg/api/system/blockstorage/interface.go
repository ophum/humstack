package blockstorage

import "github.com/gin-gonic/gin"

type BlockStorageHandlerInterface interface {
	FindAll(ctx *gin.Context)
	Find(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type BlockStorageHandler struct {
	router *gin.RouterGroup
	bshi   BlockStorageHandlerInterface
}

const (
	basePath = "namespaces/:namespace_id/blockstorages"
)

func NewBlockStorageHandler(router *gin.RouterGroup, bshi BlockStorageHandlerInterface) *BlockStorageHandler {
	return &BlockStorageHandler{
		router: router,
		bshi:   bshi,
	}
}

func (h *BlockStorageHandler) RegisterHandlers() {
	bs := h.router.Group(basePath)
	{
		bs.GET("", h.bshi.FindAll)
		bs.GET("/:block_storage_id", h.bshi.Find)
		bs.POST("", h.bshi.Create)
		bs.PUT("/:block_storage_id", h.bshi.Update)
		bs.DELETE("/:block_storage_id", h.bshi.Delete)
	}
}
