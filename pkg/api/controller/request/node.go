package request

import "github.com/ophum/humstack/v1/pkg/api/entity"

type NodeCreateRequest struct {
	Name        string                `json:"name"`
	Annotations map[string]string     `json:"annotations"`
	Hostname    string                `json:"hostname"`
	Agents      []*entity.AgentDaemon `json:"agents"`
}

type NodeUpdateStatusRequest struct {
	Status entity.NodeStatus `json:"status"`
}
