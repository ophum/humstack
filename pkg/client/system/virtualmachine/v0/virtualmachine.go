package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/system"
)

type VirtualMachineClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type VirtualMachineResponse struct {
	Code  int32       `json: "code"`
	Error interface{} `json:"error"`
	Data  struct {
		VirtualMachine system.VirtualMachine `json:"virtualMachine"`
	} `json:"data"`
}

type VirtualMachineListResponse struct {
	Code  int32       `json: "code"`
	Error interface{} `json:"error"`
	Data  struct {
		VirtualMachineList []*system.VirtualMachine `json:"virtualMachines"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/namespaces/%s/virtualmachines"
)

func NewVirtualMachineClient(scheme, apiServerAddress string, apiServerPort int32) *VirtualMachineClient {
	return &VirtualMachineClient{
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

func (c *VirtualMachineClient) Get(namespaceID, virtualMachineID string) (*system.VirtualMachine, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(namespaceID, virtualMachineID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	vmRes := VirtualMachineResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	return &vmRes.Data.VirtualMachine, nil
}

func (c *VirtualMachineClient) List(namespaceID string) ([]*system.VirtualMachine, error) {
	res, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(namespaceID, ""))
	if err != nil {
		return nil, err
	}
	body := res.Body()

	vmListRes := VirtualMachineListResponse{}
	err = json.Unmarshal(body, &vmListRes)
	if err != nil {
		return nil, err
	}

	return vmListRes.Data.VirtualMachineList, nil
}

func (c *VirtualMachineClient) Create(vm *system.VirtualMachine) (*system.VirtualMachine, error) {
	body, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(vm.Namespace, ""))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	vmRes := VirtualMachineResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error")
	}

	return &vmRes.Data.VirtualMachine, nil
}

func (c *VirtualMachineClient) Update(vm *system.VirtualMachine) (*system.VirtualMachine, error) {
	body, err := json.Marshal(vm)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(vm.Namespace, vm.ID))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	vmRes := VirtualMachineResponse{}
	err = json.Unmarshal(body, &vmRes)
	if err != nil {
		return nil, err
	}

	return &vmRes.Data.VirtualMachine, nil
}

func (c *VirtualMachineClient) Delete(namespaceID, virtualMachineID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(namespaceID, virtualMachineID))
	if err != nil {
		return err
	}
	return err

}

func (c *VirtualMachineClient) getPath(namespaceID, virtualMachineID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, namespaceID),
			virtualMachineID,
		))
}
