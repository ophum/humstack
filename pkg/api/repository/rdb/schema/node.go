package schema

import (
	"time"

	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type Node struct {
	Name        string            `gorm:"primaryKey"`
	Annotations []*NodeAnnotation `gorm:"foreignKey:NodeName"`

	Hostname string

	Status    entity.NodeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToEntityNode(v *Node) *entity.Node {
	if v == nil {
		return nil
	}
	return &entity.Node{
		Name:        v.Name,
		Annotations: ToMapNodeAnnotations(v.Annotations),
		Hostname:    v.Hostname,
		Status:      v.Status,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
	}
}

func ToEntityNodes(v []*Node) []*entity.Node {
	l := []*entity.Node{}
	for _, vv := range v {
		if vv == nil {
			continue
		}
		l = append(l, ToEntityNode(vv))
	}
	return l
}
