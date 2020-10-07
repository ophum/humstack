package v0

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/api/meta"
	"github.com/ophum/humstack/pkg/api/system"
)

type BlockStorageClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

type BlockStorageResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		BlockStorage system.BlockStorage `json:"blockStorage"`
	} `json:"data"`
}

type BlockStorageListResponse struct {
	Code  int32       `json:"code"`
	Error interface{} `json:"error"`
	Data  struct {
		BlockStorageList []*system.BlockStorage `json:"blockStorages"`
	} `json:"data"`
}

const (
	basePathFormat = "api/v0/groups/%s/namespaces/%s/blockstorages"
)

func NewBlockStorageClient(scheme, apiServerAddress string, apiServerPort int32) *BlockStorageClient {
	return &BlockStorageClient{
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

func (c *BlockStorageClient) Get(groupID, namespaceID, blockStorageID string) (*system.BlockStorage, error) {
	resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, blockStorageID))
	if err != nil {
		return nil, err
	}
	body := resp.Body()

	bsRes := BlockStorageResponse{}
	err = json.Unmarshal(body, &bsRes)
	if err != nil {
		return nil, err
	}

	return &bsRes.Data.BlockStorage, nil
}

func (c *BlockStorageClient) List(groupID, namespaceID string) ([]*system.BlockStorage, error) {
	res, err := c.client.R().SetHeaders(c.headers).Get(c.getPath(groupID, namespaceID, ""))
	if err != nil {
		return nil, err
	}
	body := res.Body()

	bsListRes := BlockStorageListResponse{}
	err = json.Unmarshal(body, &bsListRes)
	if err != nil {
		return nil, err
	}

	return bsListRes.Data.BlockStorageList, nil
}

func (c *BlockStorageClient) Create(blockstorage *system.BlockStorage) (*system.BlockStorage, error) {
	body, err := json.Marshal(blockstorage)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Post(c.getPath(blockstorage.Group, blockstorage.Namespace, ""))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	bsRes := BlockStorageResponse{}
	err = json.Unmarshal(body, &bsRes)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error: %+v", bsRes)
	}

	return &bsRes.Data.BlockStorage, nil
}

func (c *BlockStorageClient) Update(blockstorage *system.BlockStorage) (*system.BlockStorage, error) {
	body, err := json.Marshal(blockstorage)
	if err != nil {
		return nil, err
	}

	res, err := c.client.R().SetHeaders(c.headers).SetBody(body).Put(c.getPath(blockstorage.Group, blockstorage.Namespace, blockstorage.ID))
	if err != nil {
		return nil, err
	}
	body = res.Body()

	bsRes := BlockStorageResponse{}
	err = json.Unmarshal(body, &bsRes)
	if err != nil {
		return nil, err
	}

	return &bsRes.Data.BlockStorage, nil
}

func (c *BlockStorageClient) Delete(groupID, namespaceID, blockStorageID string) error {
	_, err := c.client.R().SetHeaders(c.headers).Delete(c.getPath(groupID, namespaceID, blockStorageID))
	if err != nil {
		return err
	}
	return err

}

func (c *BlockStorageClient) DeleteState(groupID, namespaceID, blockStorageID string) error {
	bs, err := c.Get(groupID, namespaceID, blockStorageID)
	if err != nil {
		return err
	}

	bs.DeleteState = meta.DeleteStateDelete

	_, err = c.Update(bs)
	return err
}

func (c *BlockStorageClient) getPath(groupID, namespaceID, blockStorageID string) string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, groupID, namespaceID),
			blockStorageID,
		))
}
