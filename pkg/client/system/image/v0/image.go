package v0

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/system"
)

type ImageClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type ImageResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		Image system.Image `json:"image"`
	} `json:"data"`
}

type ImageListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		ImageList []*system.Image `json:"images"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/images"
)

func NewImageClient(scheme, apiServerAddress string, apiServerPort int32) *ImageClient {
	return &ImageClient{
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

func (c *ImageClient) Get(groupID, imageID string) (*system.Image, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, imageID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := ImageResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Image, nil
}

func (c *ImageClient) List(groupID string) ([]*system.Image, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, ""))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	nodeResp := ImageListResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}
	return nodeResp.Data.ImageList, nil
}

func (c *ImageClient) Create(image *system.Image) (*system.Image, error) {
	body, err := json.Marshal(image)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))

	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(image.Group, ""))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := ImageResponse{}
	err = json.Unmarshal(body, &nodeResp)
	log.Println(err)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %+v", nodeResp.Error)
	}

	return &nodeResp.Data.Image, nil
}

func (c *ImageClient) Update(image *system.Image) (*system.Image, error) {
	body, err := json.Marshal(image)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	resp, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(image.Group, image.ID))
	if err != nil {
		return nil, err
	}
	body = resp.Body()

	nodeResp := ImageResponse{}
	err = json.Unmarshal(body, &nodeResp)
	if err != nil {
		return nil, err
	}

	return &nodeResp.Data.Image, nil
}

func (c *ImageClient) Delete(groupID, imageID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, imageID))
	if err != nil {
		return err
	}

	return nil
}

func (c *ImageClient) Download(groupID, imageID, tag string) (io.ReadCloser, int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/tags/%s/download", c.getPath(groupID, imageID), tag))
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, 0, fmt.Errorf("not found")
	}

	return resp.Body, int(resp.ContentLength), nil
}

func (c *ImageClient) getPath(groupID, imageID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d",
				c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID),
			imageID))
}
