package meta

import (
	"github.com/gin-gonic/gin"
)

func ResponseJSON(ctx *gin.Context, code int, err error, data interface{}) {
	if err != nil {
		ctx.JSON(code, gin.H{
			"code":  code,
			"error": err.Error(),
			"data":  data,
		})
		return
	}
	ctx.JSON(code, gin.H{
		"code":  code,
		"error": nil,
		"data":  data,
	})
}
