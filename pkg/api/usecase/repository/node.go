package repository

import "github.com/ophum/humstack/v1/pkg/api/entity"

type INodeRepository interface {
	Get(id string) (*entity.Node, error)
	List() ([]*entity.Node, error)
	Create(*entity.Node) error
	UpdateStatus(string, entity.NodeStatus) error
}
