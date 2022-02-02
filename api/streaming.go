package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
	"net/http"
)

func RegisterStreamingEndpoints(router *gin.Engine) {
	router.GET("/streaming/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		source, err := utils.StreamingRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, source)
	})
}
