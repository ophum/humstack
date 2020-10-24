package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/system"
)

type ImageEntityClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type ImageEntityResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ImageEntity system.ImageEntity `json:"imageentity"`
	} `json:"data"`
}

type ImageEntityListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ImageEntityList []*system.ImageEntity `json:"imageentities"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/imageentities"
)

func NewImageEntityClient(scheme, apiServerAddress string, apiServerPort int32) *ImageEntityClient {
	return &ImageEntityClient{
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

func (c *ImageEntityClient) Get(groupID, imageEntityID string) (*system.ImageEntity, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, imageEntityID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := ImageEntityResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.ImageEntity, nil
}

func (c *ImageEntityClient) List(groupID string) ([]*system.ImageEntity, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, ""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := ImageEntityListResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}
	return nodeResp.Data.ImageEntityList, nil
}

func (c *ImageEntityClient) Create(imageEntity *system.ImageEntity) (*system.ImageEntity, error) {
	body, err := json.Marshal(imageEntity)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(imageEntity.Group, ""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := ImageEntityResponse{}
	err = json.Unmarshal(body, &nodeResp)
	log.Println(err)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", nodeResp.Error)
	}

	return &nodeResp.Data.ImageEntity, nil
}

func (c *ImageEntityClient) Update(imageEntity *system.ImageEntity) (*system.ImageEntity, error) {
	body, err := json.Marshal(imageEntity)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(imageEntity.Group, imageEntity.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := ImageEntityResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.ImageEntity, nil
}

func (c *ImageEntityClient) Delete(groupID, imageEntityID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, imageEntityID))
	if err != nil {
		return err
	}

	return nil
}

func (c *ImageEntityClient) getPath(groupID, imageEntityID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID),
			imageEntityID))
}
