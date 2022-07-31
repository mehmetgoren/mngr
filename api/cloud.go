package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/reps"
	"net/http"
	"strconv"
)

func RegisterCloudEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/telegram", func(ctx *gin.Context) {
		cloudRep := rb.CloudRep
		viewModel := &models.TelegramViewModel{}
		viewModel.Enabled = cloudRep.IsTelegramIntegrationEnabled()
		bot, err := cloudRep.GetTelegramBot()
		if err != nil {
			return
		}
		viewModel.Bot = bot
		users, err := cloudRep.GetTelegramUsers()
		viewModel.Users = users

		ctx.JSON(http.StatusOK, viewModel)
	})
	router.POST("/telegram", func(ctx *gin.Context) {
		var model models.TelegramViewModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		cloudRep := rb.CloudRep
		_, err := cloudRep.SetTelegramIntegrationEnabled(model.Enabled)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = cloudRep.SaveTelegramBot(model.Bot)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, true)
	})
	router.DELETE("/telegramuser/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = rb.CloudRep.RemoveTelegramUserById(idInt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.GET("/gdrive", func(ctx *gin.Context) {
		cloudRep := rb.CloudRep
		viewModel := &models.GdriveViewModel{}
		viewModel.AuthCode, _ = cloudRep.GetGdriveAuthCode()
		viewModel.CredentialsJson, _ = cloudRep.GetGdriveCredentials()
		viewModel.Enabled = cloudRep.IsGdriveIntegrationEnabled()
		viewModel.TokenJson, _ = cloudRep.GetGdriveToken()
		viewModel.URL, _ = cloudRep.GetGdriveUrl()

		ctx.JSON(http.StatusOK, viewModel)
	})
	router.POST("/gdrive", func(ctx *gin.Context) {
		var model models.GdriveViewModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cloudRep := rb.CloudRep
		cloudRep.SaveGdriveAuthCode(model.AuthCode)
		cloudRep.SaveGdriveCredentials(model.CredentialsJson)
		cloudRep.SaveGdriveIntegrationEnabled(model.Enabled)
		cloudRep.SaveGdriveToken(model.TokenJson)
		cloudRep.SaveGdriveUrl(model.URL)
	})

	router.POST("/resetgdrivetokenandurl", func(ctx *gin.Context) {
		cloudRep := rb.CloudRep
		cloudRep.SaveGdriveToken("")
		cloudRep.SaveGdriveUrl("")
		ctx.JSON(http.StatusOK, true)
	})
}
