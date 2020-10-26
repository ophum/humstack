package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/core"
)

type ExternalIPClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type ExternalIPResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ExternalIP core.ExternalIP `json:"externalip"`
	} `json:"data"`
}

type ExternalIPListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ExternalIPList []*core.ExternalIP `json:"externalips"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/externalips"
)

func NewExternalIPClient(scheme, apiServerAddress string, apiServerPort int32) *ExternalIPClient {
	return &ExternalIPClient{
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

func (c *ExternalIPClient) Get(eipID string) (*core.ExternalIP, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(eipID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	eipResp := ExternalIPResponse{}
	err = json.Unmarshal(body, &eipResp)
	if err != nil {
		return nil, err
	}

	return &eipResp.Data.ExternalIP, nil
}

func (c *ExternalIPClient) List() ([]*core.ExternalIP, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	eipResp := ExternalIPListResponse{}
	err = json.Unmarshal(body, &eipResp)
	if err != nil {
		return nil, err
	}
	return eipResp.Data.ExternalIPList, nil
}

func (c *ExternalIPClient) Create(eip *core.ExternalIP) (*core.ExternalIP, error) {
	body, err := json.Marshal(eip)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	eipResp := ExternalIPResponse{}
	err = json.Unmarshal(body, &eipResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", eipResp.Error)
	}

	return &eipResp.Data.ExternalIP, nil
}

func (c *ExternalIPClient) Update(eip *core.ExternalIP) (*core.ExternalIP, error) {
	body, err := json.Marshal(eip)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(eip.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	eipResp := ExternalIPResponse{}
	err = json.Unmarshal(body, &eipResp)
	if err != nil {
		return nil, err
	}

	return &eipResp.Data.ExternalIP, nil
}

func (c *ExternalIPClient) Delete(eipID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(eipID))
	if err != nil {
		return err
	}

	return nil
}

func (c *ExternalIPClient) getPath(path string) string {
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
