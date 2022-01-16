package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

var connConfig = utils.CreateRedisConnection(utils.MAIN)

var configRep = reps.ConfigRepository{Connection: connConfig}

func RegisterMlConfigEndpoints(router *gin.Engine) {
	router.GET("/mlconfig", func(ctx *gin.Context) {
		config, _ := configRep.GetMlConfig()
		ctx.JSON(http.StatusOK, config)
	})

	router.POST("/mlconfig", func(ctx *gin.Context) {
		var config models.MlConfig
		ctx.BindJSON(&config)
		configRep.SaveMlConfig(&config)
		ctx.JSON(http.StatusOK, config)
	})

	router.GET("/restoremlconfig", func(ctx *gin.Context) {
		config, _ := configRep.RestoreMlConfig()
		ctx.JSON(http.StatusOK, config)
	})
}
