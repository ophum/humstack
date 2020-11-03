package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
)

type NamespaceClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type NamespaceResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Namespace core.Namespace `json:"namespace"`
	} `json:"data"`
}

type NamespaceListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		NamespaceList []*core.Namespace `json:"namespaces"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/namespaces"
)

func NewNamespaceClient(scheme, apiServerAddress string, apiServerPort int32) *NamespaceClient {
	return &NamespaceClient{
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

func (c *NamespaceClient) Get(groupID, namespaceID string) (*core.Namespace, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	namespaceResp := NamespaceResponse{}
	err = json.Unmarshal(body, &namespaceResp)
	if err != nil {
		return nil, err
	}

	return &namespaceResp.Data.Namespace, nil
}

func (c *NamespaceClient) List(groupID string) ([]*core.Namespace, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, ""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	namespaceResp := NamespaceListResponse{}
	err = json.Unmarshal(body, &namespaceResp)
	if err != nil {
		return nil, err
	}
	return namespaceResp.Data.NamespaceList, nil
}

func (c *NamespaceClient) Create(namespace *core.Namespace) (*core.Namespace, error) {
	body, err := json.Marshal(namespace)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(namespace.Group, ""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	namespaceResp := NamespaceResponse{}
	err = json.Unmarshal(body, &namespaceResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", namespaceResp.Error)
	}

	return &namespaceResp.Data.Namespace, nil
}

func (c *NamespaceClient) Update(namespace *core.Namespace) (*core.Namespace, error) {
	body, err := json.Marshal(namespace)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(namespace.Group, namespace.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	namespaceResp := NamespaceResponse{}
	err = json.Unmarshal(body, &namespaceResp)
	if err != nil {
		return nil, err
	}

	return &namespaceResp.Data.Namespace, nil
}

func (c *NamespaceClient) Delete(groupID, namespaceID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, namespaceID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NamespaceClient) DeleteState(groupID, namespaceID string) error {
	ns, err := c.Get(groupID, namespaceID)
	if err != nil {
		return err
	}

	ns.DeleteState = meta.DeleteStateDelete

	_, err = c.Update(ns)
	return err
}

func (c *NamespaceClient) getPath(groupID, path string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress,
				c.apiServerPort,
			),
			fmt.Sprintf(basePathFormat, groupID),
			path))
}
