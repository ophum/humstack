package controller

import (
	"log"
	"net/http"

	"github.com/ophum/humstack/v1/pkg/api/controller/request"
	"github.com/ophum/humstack/v1/pkg/api/controller/response"
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/usecase"
)

type INodeController interface {
	Get(Context)
	List(Context)
	Create(Context)
	UpdateStatus(Context)
}

const (
	paramNodeID = "node_id"
)

var _ INodeController = &NodeController{}

type NodeController struct {
	nodeUsecase usecase.INodeUsecase
}

func NewNodeController(nodeUsecase usecase.INodeUsecase) *NodeController {
	return &NodeController{nodeUsecase}
}

func (c *NodeController) Get(ctx Context) {
	id := ctx.Param(paramNodeID)
	node, err := c.nodeUsecase.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.NodeOneResponse{
		Node: node,
	})
}

func (c *NodeController) List(ctx Context) {
	nodes, err := c.nodeUsecase.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.NodeManyResponse{
		Nodes: nodes,
	})
}

func (c *NodeController) Create(ctx Context) {
	var req request.NodeCreateRequest
	if err := ctx.Bind(&req); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	node, err := c.nodeUsecase.Create(&entity.Node{
		Name:        req.Name,
		Annotations: req.Annotations,
		Hostname:    req.Hostname,
		Agents:      req.Agents,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusCreated, response.NodeOneResponse{
		Node: node,
	})
}

func (c *NodeController) UpdateStatus(ctx Context) {
	var req request.NodeUpdateStatusRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	id := ctx.Param(paramNodeID)
	if err := c.nodeUsecase.UpdateStatus(id, req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
