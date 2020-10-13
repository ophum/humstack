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
	basePathFormat = "api/v0/watches"
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

func (c *WatchClient) Watch(f func(before interface{}, after interface{})) error {
	log.Println("start")
	client := sse.NewClient(c.getPath())
	//resp, err := http.Get(c.getPath())
	//resp, err := c.client.R().SetHeaders(c.headers).Get(c.getPath())
	//if err != nil {
	//	return err
	//}
	//r := resp.RawBody()
	//r := resp.Body

	client.Subscribe("data", func(msg *sse.Event) {
		log.Println(string(msg.Data))
		var noticeData leveldb.NoticeData
		json.Unmarshal(msg.Data, &noticeData)
		log.Println(noticeData)

		f(noticeData.Before, noticeData.After)
	})

	return nil
}

func (c *WatchClient) getPath() string {
	return fmt.Sprintf("%s://%s",
		c.scheme,
		filepath.Join(
			fmt.Sprintf("%s:%d", c.apiServerAddress, c.apiServerPort),
			fmt.Sprintf(basePathFormat),
		))

}
