package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/core"
)

type ExternalIPPoolClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type ExternalIPPoolResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ExternalIPPool core.ExternalIPPool `json:"externalippool"`
	} `json:"data"`
}

type ExternalIPPoolListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ExternalIPPoolList []*core.ExternalIPPool `json:"externalippools"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/externalippools"
)

func NewExternalIPPoolClient(scheme, apiServerAddress string, apiServerPort int32) *ExternalIPPoolClient {
	return &ExternalIPPoolClient{
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

func (c *ExternalIPPoolClient) Get(eippoolID string) (*core.ExternalIPPool, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(eippoolID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	eippoolResp := ExternalIPPoolResponse{}
	err = json.Unmarshal(body, &eippoolResp)
	if err != nil {
		return nil, err
	}

	return &eippoolResp.Data.ExternalIPPool, nil
}

func (c *ExternalIPPoolClient) List() ([]*core.ExternalIPPool, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	eippoolResp := ExternalIPPoolListResponse{}
	err = json.Unmarshal(body, &eippoolResp)
	if err != nil {
		return nil, err
	}
	return eippoolResp.Data.ExternalIPPoolList, nil
}

func (c *ExternalIPPoolClient) Create(eippool *core.ExternalIPPool) (*core.ExternalIPPool, error) {
	body, err := json.Marshal(eippool)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	eippoolResp := ExternalIPPoolResponse{}
	err = json.Unmarshal(body, &eippoolResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", eippoolResp.Error)
	}

	return &eippoolResp.Data.ExternalIPPool, nil
}

func (c *ExternalIPPoolClient) Update(eippool *core.ExternalIPPool) (*core.ExternalIPPool, error) {
	body, err := json.Marshal(eippool)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(eippool.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	eippoolResp := ExternalIPPoolResponse{}
	err = json.Unmarshal(body, &eippoolResp)
	if err != nil {
		return nil, err
	}

	return &eippoolResp.Data.ExternalIPPool, nil
}

func (c *ExternalIPPoolClient) Delete(eippoolID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(eippoolID))
	if err != nil {
		return err
	}

	return nil
}

func (c *ExternalIPPoolClient) getPath(path string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress,
				c.apiServerPort,
			),
			basePathFormat,
			path))
}
