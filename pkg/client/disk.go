package client

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/v1/pkg/api/controller/request"
	"github.com/ophum/humstack/v1/pkg/api/controller/response"
	"github.com/ophum/humstack/v1/pkg/api/entity"
)

type IDiskClient interface {
	List(context.Context) ([]*entity.Disk, error)
	Create(context.Context, *entity.Disk) (*entity.Disk, error)
	UpdateStatus(context.Context, string, entity.DiskStatus) error
}

var _ IDiskClient = &DiskClient{}

type DiskClient struct {
	apiEndpoint url.URL
}

var (
	headers = map[string]string{
		"Content-Type": "application/json",
	}
)

func NewDiskClient(apiEndpoint url.URL) *DiskClient {
	return &DiskClient{apiEndpoint: apiEndpoint}
}

func (c *DiskClient) getURL(path string) url.URL {
	u := c.apiEndpoint
	u.Path = path
	return u
}

func (c *DiskClient) List(ctx context.Context) ([]*entity.Disk, error) {
	u := c.getURL("/api/v1/disks")
	client := resty.New()
	var res struct {
		Disks []*entity.Disk `json:"disks"`
	}
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetResult(&res).Get(u.String())
	if err != nil {
		return nil, err
	}
	return res.Disks, nil
}

func (c *DiskClient) Create(ctx context.Context, disk *entity.Disk) (*entity.Disk, error) {
	u := c.getURL("/api/v1/disks")
	client := resty.New()
	var res response.DiskOneResponse
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetBody(&request.DiskCreateRequest{
		Name:        disk.Name,
		Annotations: disk.Annotations,
		Type:        disk.Type,
		RequestSize: disk.RequestSize.String(),
		LimitSize:   disk.LimitSize.String(),
	}).SetResult(&res).Post(u.String())
	if err != nil {
		return nil, err
	}
	return res.Disk, nil
}

func (c *DiskClient) UpdateStatus(ctx context.Context, id string, status entity.DiskStatus) error {
	u := c.getURL("/api/v1/disks/" + id + "/status")
	client := resty.New()
	_, err := client.R().SetContext(ctx).SetHeaders(headers).SetBody(struct {
		Status entity.DiskStatus `json:"status"`
	}{
		Status: status,
	}).Patch(u.String())
	return err
}
