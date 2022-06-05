package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"net/http"
)

func RegisterServiceEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/services", func(ctx *gin.Context) {
		services, err := rb.ServiceRep.GetServices()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, services)
	})
}
