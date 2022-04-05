package api

import (
	"github.com/gin-gonic/gin"
	"mngr/eb"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"time"
)

func RegisterOdEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/ods/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		od, err := rb.OdRep.Get(id)
		if err != nil {
			//ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			ctx.JSON(http.StatusOK, nil)
			return
		}
		if len(od.Id) == 0 {
			od = nil
		}
		ctx.JSON(http.StatusOK, od)
	})
	router.POST("/ods", func(ctx *gin.Context) {
		var model models.OdModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		model.CreatedAt = utils.TimeToString(time.Now(), true)
		if _, err := rb.OdRep.Save(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mc := eb.ModelChanged{SourceId: model.Id}
		mcJson, _ := utils.SerializeJson(mc)
		eventPub := eb.DataChangedEvent{Rb: rb, ModelName: "od", ParamsJson: mcJson, Op: eb.SAVE}
		err := eventPub.Publish()
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		}
		ctx.JSON(http.StatusOK, model)
	})
}
