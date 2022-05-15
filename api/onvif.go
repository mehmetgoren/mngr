package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"net/http"
)

func RegisterOnvifEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/onvifnetwork", func(ctx *gin.Context) {
		results, err := rb.NdRep.GetAll()
		if err != nil {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusOK, results)
	})

	router.GET("/onvif/:address", func(ctx *gin.Context) {
		address := ctx.Param("address")
		results, err := rb.OvRep.Get(address)
		if err != nil {
			ctx.JSON(http.StatusOK, nil)
			return
		}
		ctx.JSON(http.StatusOK, results)
	})
}
