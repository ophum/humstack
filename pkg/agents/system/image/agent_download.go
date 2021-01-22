package image

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (a *ImageAgent) DownloadAPI(config *ImageAgentDownloadAPIConfig) error {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/api/v0/groups/:group_id/images/:image_id/tags/:tag/download", func(ctx *gin.Context) {
		groupID := ctx.Param("group_id")
		imageID := ctx.Param("image_id")
		tag := ctx.Param("tag")

		image, err := a.client.SystemV0().Image().Get(groupID, imageID)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "%v", err)
			return
		}
		imageEntityID, ok := image.Spec.EntityMap[tag]
		if !ok {
			ctx.String(http.StatusNotFound, "notfound")
			return
		}

		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename= "+imageID)

		file, err := os.Open(filepath.Join(a.localImageDirectory, groupID, imageEntityID))
		if err != nil {
			ctx.String(http.StatusInternalServerError, "%v", err)
			return
		}
		defer file.Close()

		if finfo, err := file.Stat(); err != nil {
			ctx.String(http.StatusInternalServerError, "%v", err)
			return
		} else {
			ctx.Header("Content-Length", fmt.Sprintf("%d", finfo.Size()))
		}

		_, err = io.Copy(ctx.Writer, file)
		if err != nil {
			a.logger.Error(
				"image downloader",
				zap.String("msg", err.Error()),
				zap.Time("time", time.Now()),
			)
			ctx.String(http.StatusInternalServerError, "%v", err)
		}
	})

	if err := r.Run(fmt.Sprintf("%s:%d", config.AdvertiseAddress, config.ListenPort)); err != nil {
		return err
	}

	return nil
}
