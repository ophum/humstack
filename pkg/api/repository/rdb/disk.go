package rdb

import (
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/repository/rdb/schema"
	"github.com/ophum/humstack/v1/pkg/api/usecase/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repository.IDiskRepository = &DiskRepository{}

type DiskRepository struct {
	db *gorm.DB
}

func NewDiskRepository(db *gorm.DB) *DiskRepository {
	return &DiskRepository{db}
}

func (r *DiskRepository) Get(name string) (*entity.Disk, error) {
	var disk schema.Disk
	if err := r.db.Preload(clause.Associations).Where("name = ?", name).First(&disk).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return schema.ToEntityDisk(&disk), nil
}

func (r *DiskRepository) List() ([]*entity.Disk, error) {
	var disks []*schema.Disk
	if err := r.db.Preload(clause.Associations).Find(&disks).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return schema.ToEntityDisks(disks), nil
}

func (r *DiskRepository) Create(disk *entity.Disk) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(&schema.Disk{
			Name:         disk.Name,
			Type:         disk.Type,
			RequestBytes: disk.RequestBytes,
			LimitBytes:   disk.LimitBytes,
			Status:       disk.Status,
		}).Error; err != nil {
			return err
		}

		if len(disk.Annotations) > 0 {
			annotations := []*schema.DiskAnnotation{}
			for k, v := range disk.Annotations {
				annotations = append(annotations, &schema.DiskAnnotation{
					DiskName: disk.Name,
					Key:      k,
					Value:    v,
				})
			}
			if err := r.db.Create(&annotations).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return errors.WithStack(err)
}

func (r *DiskRepository) UpdateStatus(name string, status entity.DiskStatus) error {
	err := r.db.Select("Status", "UpdatedAt").Where("name = ?", name).Updates(&schema.Disk{
		Status: status,
	}).Error
	return errors.WithStack(err)
}
