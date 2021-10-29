package schema

import (
	"time"

	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type Disk struct {
	Name        string            `gorm:"primaryKey"`
	Annotations []*DiskAnnotation `gorm:"foreignKey:DiskName"`

	Type         entity.DiskType
	RequestBytes int64
	LimitBytes   int64

	Status    entity.DiskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToEntityDisk(v *Disk) *entity.Disk {
	if v == nil {
		return nil
	}

	requestSize := entity.NewByteUnit(v.RequestBytes)
	limitSize := entity.NewByteUnit(v.LimitBytes)
	return &entity.Disk{
		Name:        v.Name,
		Annotations: ToEntityDiskAnnotations(v.Annotations),
		Type:        v.Type,
		RequestSize: *requestSize,
		LimitSize:   *limitSize,
		Status:      v.Status,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
}

func ToEntityDisks(v []*Disk) []*entity.Disk {
	l := []*entity.Disk{}
	for _, vv := range v {
		if vv == nil {
			continue
		}
		l = append(l, ToEntityDisk(vv))
	}
	return l
}
