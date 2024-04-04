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

func RegisterSmartVisionEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/smartvisions/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		sv, err := rb.SmartVisionRep.Get(id)
		if err != nil {
			//ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			ctx.JSON(http.StatusOK, nil)
			return
		}
		if len(sv.Id) == 0 {
			sv = nil
		}
		ctx.JSON(http.StatusOK, sv)
	})
	router.GET("/smartvisions", func(ctx *gin.Context) {
		svs, err := rb.SmartVisionRep.GetAll()
		if err != nil {
			svs = make([]*models.SmartVisionModel, 0)
		}
		ctx.JSON(http.StatusOK, svs)
	})
	router.POST("/smartvisions", func(ctx *gin.Context) {
		var model models.SmartVisionModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		model.CreatedAt = utils.TimeToString(time.Now(), true)
		if _, err := rb.SmartVisionRep.Save(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mc := eb.ModelChanged{SourceId: model.Id}
		mcJson, _ := utils.SerializeJson(mc)
		eventPub := eb.DataChangedEvent{Rb: rb, ModelName: "smart_vision", ParamsJson: mcJson, Op: eb.SAVE}
		err := eventPub.Publish()
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		}
		ctx.JSON(http.StatusOK, model)
	})
}
