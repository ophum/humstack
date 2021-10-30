package rdb

import (
	"encoding/json"

	"github.com/ophum/humstack/v1/pkg/api/entity"
	"github.com/ophum/humstack/v1/pkg/api/repository/rdb/schema"
	"github.com/ophum/humstack/v1/pkg/api/usecase/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repository.INodeRepository = &NodeRepository{}

type NodeRepository struct {
	db *gorm.DB
}

func NewNodeRepository(db *gorm.DB) *NodeRepository {
	return &NodeRepository{db}
}

func (r *NodeRepository) Get(name string) (*entity.Node, error) {
	var node schema.Node
	if err := r.db.Preload(clause.Associations).Where("name = ?", name).First(&node).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	var agents []*schema.NodeAgentDaemon
	if err := r.db.Where("node_name = ?", name).Find(&agents).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	entityAgents := schema.ToEntityAgentDaemons(agents)
	entityNode := schema.ToEntityNode(&node)
	entityNode.Agents = entityAgents
	return entityNode, nil
}

func (r *NodeRepository) List() ([]*entity.Node, error) {
	var nodes []*schema.Node
	if err := r.db.Preload(clause.Associations).Find(&nodes).Error; err != nil {
		return nil, errors.WithStack(err)
	}

	entityNodes := schema.ToEntityNodes(nodes)
	for i, node := range nodes {
		var agents []*schema.NodeAgentDaemon
		if err := r.db.Where("node_name = ?", node.Name).Find(&agents).Error; err != nil {
			return nil, errors.WithStack(err)
		}
		entityNodes[i].Agents = schema.ToEntityAgentDaemons(agents)
	}
	return entityNodes, nil
}

func (r *NodeRepository) Create(node *entity.Node) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&schema.Node{
			Name:     node.Name,
			Hostname: node.Hostname,
			Status:   node.Status,
		}).Error; err != nil {
			return err
		}

		if len(node.Annotations) > 0 {
			annotations := []*schema.NodeAnnotation{}
			for k, v := range node.Annotations {
				annotations = append(annotations, &schema.NodeAnnotation{
					NodeName: node.Name,
					Key:      k,
					Value:    v,
				})
			}
			if err := tx.Create(&annotations).Error; err != nil {
				return err
			}
		}
		if len(node.Agents) > 0 {
			agents := []*schema.NodeAgentDaemon{}
			for _, agent := range node.Agents {
				args, err := json.Marshal(agent.Args)
				if err != nil {
					return err
				}
				envs, err := json.Marshal(agent.Envs)
				if err != nil {
					return err
				}
				agents = append(agents, &schema.NodeAgentDaemon{
					NodeName:      node.Name,
					Name:          agent.Name,
					Command:       agent.Command,
					Args:          string(args),
					Envs:          string(envs),
					RestartPolicy: string(agent.RestartPolicy),
				})
			}
			if err := tx.Create(&agents).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return errors.WithStack(err)
}

func (r *NodeRepository) UpdateStatus(name string, status entity.NodeStatus) error {
	err := r.db.Select("Status", "UpdatedAt").Where("name = ?", name).Updates(&schema.Node{
		Status: status,
	}).Error
	return errors.WithStack(err)
}
