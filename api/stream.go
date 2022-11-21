package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"net/http"
)

func RegisterStreamEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/stream", func(c *gin.Context) {
		modelList, err := rb.StreamRep.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, modelList)
	})
	router.GET("/stream/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		stream, err := rb.StreamRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, stream)
	})
}
