package client

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/v1/pkg/api/controller/request"
	"github.com/ophum/humstack/v1/pkg/api/controller/response"
	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type INodeClient interface {
	Get(ctx context.Context, hostname string) (*entity.Node, error)
	List(context.Context) ([]*entity.Node, error)
	Create(context.Context, *entity.Node) (*entity.Node, error)
	UpdateStatus(context.Context, string, entity.NodeStatus) error
}

var _ INodeClient = &NodeClient{}

type NodeClient struct {
	apiEndpoint url.URL
}

func NewNodeClient(apiEndpoint url.URL) *NodeClient {
	return &NodeClient{apiEndpoint: apiEndpoint}
}

func (c *NodeClient) getURL(path string) url.URL {
	u := c.apiEndpoint
	u.Path = path
	return u
}

func (c *NodeClient) Get(ctx context.Context, hostname string) (*entity.Node, error) {
	u := c.getURL("/api/v1/nodes/" + hostname)
	client := resty.New()
	var res response.NodeOneResponse
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetResult(&res).Get(u.String())
	return res.Node, err
}

func (c *NodeClient) List(ctx context.Context) ([]*entity.Node, error) {
	u := c.getURL("/api/v1/nodes")
	client := resty.New()
	var res struct {
		Nodes []*entity.Node `json:"nodes"`
	}
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetResult(&res).Get(u.String())
	if err != nil {
		return nil, err
	}
	return res.Nodes, nil
}

func (c *NodeClient) Create(ctx context.Context, node *entity.Node) (*entity.Node, error) {
	u := c.getURL("/api/v1/nodes")
	client := resty.New()
	var res response.NodeOneResponse
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetBody(&request.NodeCreateRequest{
		Name:        node.Name,
		Annotations: node.Annotations,
		Hostname:    node.Hostname,
		Agents:      node.Agents,
	}).SetResult(&res).Post(u.String())
	if err != nil {
		return nil, err
	}
	return res.Node, nil
}

func (c *NodeClient) UpdateStatus(ctx context.Context, id string, status entity.NodeStatus) error {
	u := c.getURL("/api/v1/nodes/" + id + "/status")
	client := resty.New()
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetBody(struct {
		Status entity.NodeStatus `json:"status"`
	}{
		Status: status,
	}).Patch(u.String())
	return err
}
