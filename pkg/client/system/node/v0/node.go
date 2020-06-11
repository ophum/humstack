package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/system"
)

type NodeClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type NodeResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Node system.Node `json:"node"`
	} `json:"data"`
}

type NodeListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		NodeList []*system.Node `json:"nodes"`
	} `json:"data"`
}

const (
	basePath = "api/v0/nodes"
)

func NewNodeClient(scheme, apiServerAddress string, apiServerPort int32) *NodeClient {
	return &NodeClient{
		scheme:           scheme,
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,
		client:           resty.New(),
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accepted":     "application/json",
		},
	}
}

func (c *NodeClient) Get(nodeID string) (*system.Node, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(nodeID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NodeResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Node, nil
}

func (c *NodeClient) List() ([]*system.Node, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NodeListResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}
	return nodeResp.Data.NodeList, nil
}

func (c *NodeClient) Create(node *system.Node) (*system.Node, error) {
	body, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NodeResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", nodeResp.Error)
	}

	return &nodeResp.Data.Node, nil
}

func (c *NodeClient) Update(node *system.Node) (*system.Node, error) {
	body, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(node.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NodeResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Node, nil
}

func (c *NodeClient) Delete(nodeID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(nodeID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NodeClient) getPath(path string) string {
	return fmt.Sprintf("%s://%s", c.scheme, filepath.Join(fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort), basePath, path))
}
