package request

import (
	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type DiskCreateRequest struct {
	Name        string            `json:"name"`
	Annotations map[string]string `json:"annotations"`
	Type        entity.DiskType   `json:"disk_type"`
	RequestSize string            `json:"request_size"`
	LimitSize   string            `json:"limit_size"`
}

func (r DiskCreateRequest) ParseRequestSize() (*entity.ByteUnit, error) {
	return entity.ParseByteUnit(r.RequestSize)
}

func (r DiskCreateRequest) ParseLimitSize() (*entity.ByteUnit, error) {
	return entity.ParseByteUnit(r.LimitSize)
}

type DiskUpdateStatusRequest struct {
	Status entity.DiskStatus `json:"status"`
}
