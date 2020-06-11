package v0

import (
	"encoding/json"
	"fmt"
	"log"
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
	basePathFormat = "api/v0/namespaces/%s/networks"
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

func (c *NetworkClient) Get(namespaceID, networkID string) (*system.Network, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(namespaceID, networkID))
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

func (c *NetworkClient) List(namespaceID string) ([]*system.Network, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(namespaceID, ""))
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
	log.Println(string(body))

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(network.Namespace, ""))
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

	fmt.Println(string(body))
	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(network.Namespace, network.ID))
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

func (c *NetworkClient) Delete(namespaceID, networkID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(namespaceID, networkID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NetworkClient) getPath(namespaceID, networkID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, namespaceID),
			networkID))
}
