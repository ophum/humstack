package meta

import "github.com/gin-gonic/gin"

func ResponseJSON(ctx *gin.Context, code int, err error, data interface{}) {
	ctx.JSON(code, gin.H{
		"code":  code,
		"error": err,
		"data":  data,
	})
}
