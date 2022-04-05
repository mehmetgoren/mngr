package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/reps"
	"net/http"
)

func RegisterConfigEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/config", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()
		ctx.JSON(http.StatusOK, config)
	})

	router.POST("/config", func(ctx *gin.Context) {
		var config models.Config
		ctx.BindJSON(&config)
		rb.ConfigRep.SaveConfig(&config)
		ctx.JSON(http.StatusOK, config)
	})

	router.GET("/restoreconfig", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.RestoreConfig()
		ctx.JSON(http.StatusOK, config)
	})
}
