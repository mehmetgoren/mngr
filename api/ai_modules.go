package api

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"mngr/dckr"
	"mngr/models"
	"mngr/reps"
	"net/http"
)

func RegisterAiModulesEndpoints(router *gin.Engine, rb *reps.RepoBucket, dockerClient *client.Client) {
	router.GET("/aimodules", func(ctx *gin.Context) {
		aiModules, _ := rb.AiModuleRep.GetAiModules()
		ctx.JSON(http.StatusOK, aiModules)
	})

	router.POST("/aimodule", func(ctx *gin.Context) {
		var aiModules *models.AiModuleModel
		ctx.BindJSON(&aiModules)
		rb.AiModuleRep.SaveAiModule(aiModules)

		restartServiceInstance(dockerClient)

		ctx.JSON(http.StatusOK, aiModules)
	})

	router.DELETE("/aimodule/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		result, err := rb.AiModuleRep.RemoveAiModuleByName(name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			restartServiceInstance(dockerClient)

			ctx.JSON(http.StatusOK, result)
		}
	})
}

func restartServiceInstance(dockerClient *client.Client) {
	dm := dckr.DockerManager{Client: dockerClient}
	dm.RestartContainer("senseai_service-instance")
}
