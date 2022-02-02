package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
	"net/http"
)

func RegisterStreamingEndpoints(router *gin.Engine) {

	router.GET("/streaming", func(c *gin.Context) {
		models, err := utils.StreamingRep.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, models)
	})

	router.GET("/streaming/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		source, err := utils.StreamingRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, source)
	})
}
