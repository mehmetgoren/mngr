package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

var connSources = utils.CreateRedisConnection(utils.SOURCES)

var sourceRep = reps.SourceRepository{Connection: connSources}

func RegisterSourceEndpoints(router *gin.Engine) {
	router.GET("/sources", func(ctx *gin.Context) {
		sources, _ := sourceRep.GetAllSources()
		ctx.JSON(http.StatusOK, sources)
	})
}
