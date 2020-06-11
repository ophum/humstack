package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/core"
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
	basePath = "api/v0/namespaces"
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

func (c *NamespaceClient) Get(namespaceID string) (*core.Namespace, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(namespaceID))
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

func (c *NamespaceClient) List() ([]*core.Namespace, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
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
	log.Println(string(body))

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
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

	fmt.Println(string(body))
	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(namespace.ID))
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

func (c *NamespaceClient) Delete(namespaceID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(namespaceID))
	if err != nil {
		return err
	}

	return nil
}

func (c *NamespaceClient) getPath(path string) string {
	return fmt.Sprintf("%s://%s", c.scheme, filepath.Join(fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort), basePath, path))
}
