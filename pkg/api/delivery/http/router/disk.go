package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/v1/pkg/api/controller"
)

type DiskRouter struct {
	r              *gin.Engine
	diskController controller.IDiskController
}

func NewDiskRouter(
	r *gin.Engine,
	diskController controller.IDiskController,
) *DiskRouter {
	return &DiskRouter{r, diskController}
}

func (r *DiskRouter) RegisterRoutes() {
	v1 := r.r.Group("/api/v1/disks")
	{
		v1.GET("", func(c *gin.Context) {
			r.diskController.List(c)
		})
		v1.GET("/:disk_id", func(c *gin.Context) {
			r.diskController.Get(c)
		})
		v1.POST("", func(c *gin.Context) {
			r.diskController.Create(c)
		})
		v1.PATCH("/:disk_id/status", func(c *gin.Context) {
			r.diskController.UpdateStatus(c)
		})
	}
}
