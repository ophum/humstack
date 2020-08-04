package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

type VirtualRouterClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type VirtualRouterResponse struct {
	Code  int32       `json: "code"`
	Error interface{} `json:"error"`
	Data  struct {
		VirtualRouter system.VirtualRouter `json:"virtualRouter"`
	} `json:"data"`
}

type VirtualRouterListResponse struct {
	Code  int32       `json: "code"`
	Error interface{} `json:"error"`
	Data  struct {
		VirtualRouterList []*system.VirtualRouter `json:"virtualRouters"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/namespaces/%s/virtualrouters"
)

func NewVirtualRouterClient(scheme, apiServerAddress string, apiServerPort int32) *VirtualRouterClient {
	return &VirtualRouterClient{
		scheme:           scheme,
		apiServerAddress: apiServerAddress,
		apiServerPort:    apiServerPort,
		client:           resty.New(),
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

func (c *VirtualRouterClient) Get(groupID, namespaceID, virtualRouterID string) (*system.VirtualRouter, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, virtualRouterID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	vmRes := VirtualRouterResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	return &vmRes.Data.VirtualRouter, nil
}

func (c *VirtualRouterClient) List(groupID, namespaceID string) ([]*system.VirtualRouter, error) {
	res, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, ""))
	if err != nil {
		return nil, err
	}
	body := res.Body()

	vmListRes := VirtualRouterListResponse{}
	err = json.Unmarshal(body, &vmListRes)
	if err != nil {
		return nil, err
	}

	return vmListRes.Data.VirtualRouterList, nil
}

func (c *VirtualRouterClient) Create(vm *system.VirtualRouter) (*system.VirtualRouter, error) {
	body, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(vm.Group, vm.Namespace, ""))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	vmRes := VirtualRouterResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error")
	}

	return &vmRes.Data.VirtualRouter, nil
}

func (c *VirtualRouterClient) Update(vm *system.VirtualRouter) (*system.VirtualRouter, error) {
	body, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(vm.Group, vm.Namespace, vm.ID))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	vmRes := VirtualRouterResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	return &vmRes.Data.VirtualRouter, nil
}

func (c *VirtualRouterClient) Delete(groupID, namespaceID, virtualRouterID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, namespaceID, virtualRouterID))
	if err != nil {
		return err
	}
	return err

}

func (c *VirtualRouterClient) DeleteState(groupID, namespaceID, virtualRouterID string) error {
	vm, err := c.Get(groupID, namespaceID, virtualRouterID)
	if err != nil {
		return err
	}

	vm.DeleteState = meta.DeleteStateDelete

	_, err = c.Update(vm)
	return err
}

func (c *VirtualRouterClient) getPath(groupID, namespaceID, virtualRouterID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID, namespaceID),
			virtualRouterID,
		))
}
