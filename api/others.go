package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"net/http"
)

func RegisterOthersEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
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
}
