package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ophum/humstack/v1/pkg/api/controller"
)

type NodeRouter struct {
	r              *gin.Engine
	nodeController controller.INodeController
}

func NewNodeRouter(
	r *gin.Engine,
	nodeController controller.INodeController,
) *NodeRouter {
	return &NodeRouter{r, nodeController}
}

func (r *NodeRouter) RegisterRoutes() {
	v1 := r.r.Group("/api/v1/nodes")
	{
		v1.GET("", func(c *gin.Context) {
			r.nodeController.List(c)
		})
		v1.GET("/:node_id", func(c *gin.Context) {
			r.nodeController.Get(c)
		})
		v1.POST("", func(c *gin.Context) {
			r.nodeController.Create(c)
		})
		v1.PATCH("/:node_id/status", func(c *gin.Context) {
			r.nodeController.UpdateStatus(c)
		})
	}
}
