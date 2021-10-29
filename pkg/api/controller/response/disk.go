package response

import "github.com/ophum/humstack/v1/pkg/api/entity"

type DiskOneResponse struct {
	Disk *entity.Disk `json:"disk"`
}

type DiskManyResponse struct {
	Disks []*entity.Disk `json:"disks"`
}
