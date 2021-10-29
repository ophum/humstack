package controller

import (
	"net/http"

	"github.com/ophum/humstack/v1/pkg/api/controller/request"
	"github.com/ophum/humstack/v1/pkg/api/controller/response"
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/usecase"
)

type Context interface {
	Bind(interface{}) error
	Param(string) string
	JSON(int, interface{})
}
type IDiskController interface {
	Get(Context)
	List(Context)
	Create(Context)
	UpdateStatus(Context)
}

var _ IDiskController = &DiskController{}

type DiskController struct {
	diskUsecase usecase.IDiskUsecase
}

func NewDiskController(diskUsecase usecase.IDiskUsecase) *DiskController {
	return &DiskController{diskUsecase}
}

func (c *DiskController) Get(ctx Context) {
	id := ctx.Param("disk_id")
	disk, err := c.diskUsecase.Get(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.DiskOneResponse{
		Disk: disk,
	})
}

func (c *DiskController) List(ctx Context) {
	disks, err := c.diskUsecase.List()
	if err != nil {
		ctx.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, response.DiskManyResponse{
		Disks: disks,
	})
}

func (c *DiskController) Create(ctx Context) {
	var req request.DiskCreateRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	disk, err := c.diskUsecase.Create(&entity.Disk{
		Name:         req.Name,
		Annotations:  req.Annotations,
		Type:         req.Type,
		RequestBytes: req.RequestBytes,
		LimitBytes:   req.LimitBytes,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.JSON(http.StatusCreated, response.DiskOneResponse{
		Disk: disk,
	})
}

func (c *DiskController) UpdateStatus(ctx Context) {
	var req request.DiskUpdateStatusRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	id := ctx.Param("disk_id")
	if err := c.diskUsecase.UpdateStatus(id, req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
