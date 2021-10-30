package request

import (
	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type DiskCreateRequest struct {
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`
	Type        entity.DiskType   `json:"disk_type"`
	RequestSize entity.ByteUnit   `json:"request_size"`
	LimitSize   entity.ByteUnit   `json:"limit_size"`
}

type DiskUpdateStatusRequest struct {
	Status entity.DiskStatus `json:"status"`
}
