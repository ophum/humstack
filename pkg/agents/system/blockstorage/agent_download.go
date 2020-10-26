package blockstorage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *BlockStorageAgent) DownloadAPI(config *BlockStorageAgentDownloadAPIConfig) error {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.GET("/api/v0/groups/:group_id/namespaces/:namespace_id/blockstorages/:bs_id/download", func(ctx *gin.Context) {
		groupID := ctx.Param("group_id")
		namespaceID := ctx.Param("namespace_id")
		bsID := ctx.Param("bs_id")

		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename= "+bsID)

		file, err := os.Open(filepath.Join(a.localBlockStorageDirectory, groupID, namespaceID, bsID))
		if err != nil {
			ctx.String(http.StatusInternalServerError, "%v", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(ctx.Writer, file)
		if err != nil {
			log.Println(err)
		}
	})

	if err := r.Run(fmt.Sprintf("%s:%d", config.AdvertiseAddress, config.ListenPort)); err != nil {
		return err
	}

	return nil
}
