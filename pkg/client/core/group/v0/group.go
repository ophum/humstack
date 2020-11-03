package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/core"
	"github.com/ophum/humstack/pkg/api/meta"
)

type GroupClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type GroupResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Group core.Group `json:"group"`
	} `json:"data"`
}

type GroupListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		GroupList []*core.Group `json:"groups"`
	} `json:"data"`
}

const (
	basePath = "api/v0/groups"
)

func NewGroupClient(scheme, apiServerAddress string, apiServerPort int32) *GroupClient {
	return &GroupClient{
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

func (c *GroupClient) Get(groupID string) (*core.Group, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	groupResp := GroupResponse{}
	err = json.Unmarshal(body, &groupResp)
	if err != nil {
		return nil, err
	}

	return &groupResp.Data.Group, nil
}

func (c *GroupClient) List() ([]*core.Group, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	groupResp := GroupListResponse{}
	err = json.Unmarshal(body, &groupResp)
	if err != nil {
		return nil, err
	}
	return groupResp.Data.GroupList, nil
}

func (c *GroupClient) Create(group *core.Group) (*core.Group, error) {
	body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	groupResp := GroupResponse{}
	err = json.Unmarshal(body, &groupResp)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", groupResp.Error)
	}

	return &groupResp.Data.Group, nil
}

func (c *GroupClient) Update(group *core.Group) (*core.Group, error) {
	body, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(group.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	groupResp := GroupResponse{}
	err = json.Unmarshal(body, &groupResp)
	if err != nil {
		return nil, err
	}

	return &groupResp.Data.Group, nil
}

func (c *GroupClient) Delete(groupID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID))
	if err != nil {
		return err
	}

	return nil
}

func (c *GroupClient) DeleteState(groupID string) error {
	group, err := c.Get(groupID)
	if err != nil {
		return err
	}

	group.DeleteState = meta.DeleteStateDelete

	_, err = c.Update(group)
	return err
}

func (c *GroupClient) getPath(path string) string {
	return fmt.Sprintf("%s://%s", c.scheme, filepath.Join(fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort), basePath, path))
}
