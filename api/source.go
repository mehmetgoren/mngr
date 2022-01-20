package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
	"net/http"
)

func RegisterSourceEndpoints(router *gin.Engine) {
	router.GET("/sources", func(ctx *gin.Context) {
		sources, _ := utils.SourceRep.GetAllSources()
		ctx.JSON(http.StatusOK, sources)
	})
}
