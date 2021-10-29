package usecase

import (
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/usecase/repository"
	"github.com/pkg/errors"
)

type IDiskUsecase interface {
	Get(string) (*entity.Disk, error)
	List() ([]*entity.Disk, error)
	Create(*entity.Disk) (*entity.Disk, error)
	UpdateStatus(string, entity.DiskStatus) error
}

var _ IDiskUsecase = &DiskUsecase{}

type DiskUsecase struct {
	diskRepo repository.IDiskRepository
}

func NewDiskUsecase(diskRepo repository.IDiskRepository) *DiskUsecase {
	return &DiskUsecase{diskRepo}
}

func (u *DiskUsecase) Get(name string) (*entity.Disk, error) {
	disk, err := u.diskRepo.Get(name)
	return disk, errors.WithStack(err)
}

func (u *DiskUsecase) List() ([]*entity.Disk, error) {
	disks, err := u.diskRepo.List()
	return disks, errors.WithStack(err)
}

func (u *DiskUsecase) Create(disk *entity.Disk) (*entity.Disk, error) {
	disk.Status = entity.DiskStatusPending
	if err := u.diskRepo.Create(disk); err != nil {
		return nil, errors.WithStack(err)
	}
	created, err := u.diskRepo.Get(disk.Name)
	return created, errors.WithStack(err)
}

func (u *DiskUsecase) UpdateStatus(name string, status entity.DiskStatus) error {
	err := u.diskRepo.UpdateStatus(name, status)
	return errors.WithStack(err)
}
