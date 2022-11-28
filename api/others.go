package api

import (
	"github.com/gin-gonic/gin"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterOthersEndpoints(router *gin.Engine, rb *reps.RepoBucket, global *models.GlobalModel) {
	router.GET("/rtsptemplates", func(ctx *gin.Context) {
		ret, err := rb.RtspTemplateRep.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, ret)
	})

	router.GET("/failedstreams", func(ctx *gin.Context) {
		ret, err := rb.FailedStreamRep.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, ret)
	})

	router.GET("/recstucks", func(ctx *gin.Context) {
		ret, err := rb.RecStuckRep.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, ret)
	})

	router.GET("/various", func(ctx *gin.Context) {
		ret, err := rb.VariousRep.Get()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, ret)
	})

	router.GET("nvidiasmi", func(ctx *gin.Context) {
		viewModel := utils.NvidiaGpuModel{}
		err := viewModel.Fetch()
		if err != nil {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusOK, &viewModel)
	})

	router.GET("isReadOnlyMode", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, global.ReadOnlyMode)
	})
}
