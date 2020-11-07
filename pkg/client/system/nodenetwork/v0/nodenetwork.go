package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

type NodeNetworkClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type NodeNetworkResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		NodeNetwork system.NodeNetwork `json:"nodenetwork"`
	} `json:"data"`
}

type NodeNetworkListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		NodeNetworkList []*system.NodeNetwork `json:"nodenetworks"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/namespaces/%s/nodenetworks"
)

func NewNodeNetworkClient(scheme, apiServerAddress string, apiServerPort int32) *NodeNetworkClient {
	return &NodeNetworkClient{
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

func (c *NodeNetworkClient) Get(groupID, namespaceID, nodenetworkID string) (*system.NodeNetwork, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, nodenetworkID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NodeNetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.NodeNetwork, nil
}

func (c *NodeNetworkClient) List(groupID, namespaceID string) ([]*system.NodeNetwork, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, ""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NodeNetworkListResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}
	return nodeResp.Data.NodeNetworkList, nil
}

func (c *NodeNetworkClient) Create(nodenetwork *system.NodeNetwork) (*system.NodeNetwork, error) {
	body, err := json.Marshal(nodenetwork)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(nodenetwork.Group, nodenetwork.Namespace, ""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NodeNetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", nodeResp.Error)
	}

	return &nodeResp.Data.NodeNetwork, nil
}

func (c *NodeNetworkClient) Update(nodenetwork *system.NodeNetwork) (*system.NodeNetwork, error) {
	body, err := json.Marshal(nodenetwork)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(nodenetwork.Group, nodenetwork.Namespace, nodenetwork.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NodeNetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.NodeNetwork, nil
}

func (c *NodeNetworkClient) Delete(groupID, namespaceID, nodenetworkID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, namespaceID, nodenetworkID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NodeNetworkClient) DeleteState(groupID, namespaceID, nodenetworkID string) error {
	net, err := c.Get(groupID, namespaceID, nodenetworkID)
	if err != nil {
		return err
	}

	net.DeleteState = meta.DeleteStateDelete

	_, err = c.Update(net)
	return err
}

func (c *NodeNetworkClient) getPath(groupID, namespaceID, nodenetworkID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID, namespaceID),
			nodenetworkID))
}
