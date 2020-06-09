package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/pkg/api/system/blockstorage"
	"github.com/ophum/humstack/pkg/store"
)

type BlockStorageHandler struct {
	blockstorage.BlockStorageHandlerInterface

	store store.Store
}

func NewBlockStorageHandler(store store.Store) *BlockStorageHandler {
	return &BlockStorageHandler{
		store: store,
	}
}

func (h *BlockStorageHandler) FindAll(ctx *gin.Context) {

}
func (h *BlockStorageHandler) Find(ctx *gin.Context) {

}
func (h *BlockStorageHandler) Create(ctx *gin.Context) {

}
func (h *BlockStorageHandler) Update(ctx *gin.Context) {

}
func (h *BlockStorageHandler) Delete(ctx *gin.Context) {

}
