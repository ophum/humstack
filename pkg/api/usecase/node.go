package usecase

import (
	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/usecase/repository"
	"github.com/pkg/errors"
)

type INodeUsecase interface {
	Get(string) (*entity.Node, error)
	List() ([]*entity.Node, error)
	Create(*entity.Node) (*entity.Node, error)
	UpdateStatus(string, entity.NodeStatus) error
}

var _ INodeUsecase = &NodeUsecase{}

type NodeUsecase struct {
	nodeRepo repository.INodeRepository
}

func NewNodeUsecase(nodeRepo repository.INodeRepository) *NodeUsecase {
	return &NodeUsecase{nodeRepo}
}

func (u *NodeUsecase) Get(name string) (*entity.Node, error) {
	node, err := u.nodeRepo.Get(name)
	return node, errors.WithStack(err)
}

func (u *NodeUsecase) List() ([]*entity.Node, error) {
	nodes, err := u.nodeRepo.List()
	return nodes, errors.WithStack(err)
}

func (u *NodeUsecase) Create(node *entity.Node) (*entity.Node, error) {
	node.Status = entity.NodeStatusNotReady
	if err := u.nodeRepo.Create(node); err != nil {
		return nil, errors.WithStack(err)
	}
	created, err := u.nodeRepo.Get(node.Name)
	return created, errors.WithStack(err)
}

func (u *NodeUsecase) UpdateStatus(name string, status entity.NodeStatus) error {
	err := u.nodeRepo.UpdateStatus(name, status)
	return errors.WithStack(err)
}
