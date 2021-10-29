package repository

import "github.com/ophum/humstack/v1/pkg/api/entity"

type IDiskRepository interface {
	Get(id string) (*entity.Disk, error)
	List() ([]*entity.Disk, error)
	Create(*entity.Disk) error
	UpdateStatus(string, entity.DiskStatus) error
}
