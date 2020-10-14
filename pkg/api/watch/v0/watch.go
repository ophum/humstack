package v0

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ophum/humstack/pkg/api/watch"
	"github.com/ophum/humstack/pkg/store/leveldb"
)

type WatchHandler struct {
	watch.WatchHandlerInterface

	notifiers map[string](chan string)
}

func NewWatchHandler(notifiers map[string](chan string)) *WatchHandler {
	return &WatchHandler{
		notifiers: notifiers,
	}
}

func (h *WatchHandler) Watch(ctx *gin.Context) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}

	idString := id.String()
	h.notifiers[idString] = make(chan string)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	apiType := ctx.DefaultQuery("apiType", "")

	w := ctx.Writer
	go func() {

		for s := range h.notifiers[idString] {
			noticeData := leveldb.NoticeData{}
			json.Unmarshal([]byte(s), &noticeData)

			if apiType == "" || apiType == string(noticeData.APIType) {

				w.Write([]byte(fmt.Sprintf("data: %s\n\n", s)))
				w.Flush()
			}

		}

	}()
	<-ctx.Done()

}
