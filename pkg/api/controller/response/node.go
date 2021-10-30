package response

import "github.com/ophum/humstack/v1/pkg/api/entity"

type NodeOneResponse struct {
	Node *entity.Node `json:"node"`
}

type NodeManyResponse struct {
	Nodes []*entity.Node `json:"nodes"`
}
