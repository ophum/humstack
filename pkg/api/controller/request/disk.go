package request

import "github.com/ophum/humstack/v1/pkg/api/entity"

type DiskCreateRequest struct {
	Name         string            `json:"name"`
	Annotations  map[string]string `json:"annotations"`
	Type         entity.DiskType   `json:"disk_type"`
	RequestBytes int               `json:"requset_bytes"`
	LimitBytes   int               `json:"limit_bytes"`
}

type DiskUpdateStatusRequest struct {
	Status entity.DiskStatus `json:"status"`
}
