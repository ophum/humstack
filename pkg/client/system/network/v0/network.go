package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/system"
)

type NetworkClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type NetworkResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Network system.Network `json:"network"`
	} `json:"data"`
}

type NetworkListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		NetworkList []*system.Network `json:"networks"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/namespaces/%s/networks"
)

func NewNetworkClient(scheme, apiServerAddress string, apiServerPort int32) *NetworkClient {
	return &NetworkClient{
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

func (c *NetworkClient) Get(groupID, namespaceID, networkID string) (*system.Network, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, networkID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Network, nil
}

func (c *NetworkClient) List(groupID, namespaceID string) ([]*system.Network, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, ""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := NetworkListResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}
	return nodeResp.Data.NetworkList, nil
}

func (c *NetworkClient) Create(network *system.Network) (*system.Network, error) {
	body, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(network.Group, network.Namespace, ""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", nodeResp.Error)
	}

	return &nodeResp.Data.Network, nil
}

func (c *NetworkClient) Update(network *system.Network) (*system.Network, error) {
	body, err := json.Marshal(network)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(network.Group, network.Namespace, network.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := NetworkResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Network, nil
}

func (c *NetworkClient) Delete(groupID, namespaceID, networkID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, namespaceID, networkID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NetworkClient) getPath(groupID, namespaceID, networkID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID, namespaceID),
			networkID))
}
