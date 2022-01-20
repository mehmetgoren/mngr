package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/utils"
	"net/http"
)

func RegisterConfigEndpoints(router *gin.Engine) {
	router.GET("/config", func(ctx *gin.Context) {
		config, _ := utils.ConfigRep.GetConfig()
		ctx.JSON(http.StatusOK, config)
	})

	router.POST("/config", func(ctx *gin.Context) {
		var config models.Config
		ctx.BindJSON(&config)
		utils.ConfigRep.SaveConfig(&config)
		ctx.JSON(http.StatusOK, config)
	})

	router.GET("/restoreconfig", func(ctx *gin.Context) {
		config, _ := utils.ConfigRep.RestoreConfig()
		ctx.JSON(http.StatusOK, config)
	})
}
