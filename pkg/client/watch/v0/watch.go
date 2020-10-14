package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/ophum/humstack/pkg/store/leveldb"
	"github.com/r3labs/sse"
)

type WatchClient struct {
	scheme           string
	apiServerAddress string
	apiServerPort    int32
	client           *resty.Client
	headers          map[string]string
}

const (
	basePathFormat = "api/v0/watches%s"
)

func NewWatchClient(scheme, apiServerAddress string, apiServerPort int32) *WatchClient {
	return &WatchClient{
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

func (c *WatchClient) Watch(apiType string, f func(before interface{}, after interface{})) error {
	log.Println("start")
	client := sse.NewClient(c.getPath(apiType))

	client.Subscribe("", func(msg *sse.Event) {
		var noticeData leveldb.NoticeData
		json.Unmarshal(msg.Data, &noticeData)

		f(noticeData.Before, noticeData.After)
	})

	return nil
}

func (c *WatchClient) getPath(apiType string) string {
	query := ""
	if apiType != "" {
		query = "?apiType=" + apiType
	}
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat, query),
		))

}
