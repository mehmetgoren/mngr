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

func RegisterSourceEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/sources", func(ctx *gin.Context) {
		sources, _ := rb.SourceRep.GetAll()
		ctx.JSON(http.StatusOK, sources)
	})
	router.GET("/sources/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		source, err := rb.SourceRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, source)
	})
	router.POST("/sources", func(ctx *gin.Context) {
		var model models.SourceModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		isNew := len(model.Id) == 0
		model.CreatedAt = utils.TimeToString(time.Now(), true)
		if _, err := rb.SourceRep.Save(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//restart the stream after edit.
		if !isNew {
			eventPub := eb.RestartStreamRequestEvent{Rb: rb, SourceModel: model}
			err := eventPub.Publish()
			if err == nil {
				ctx.Writer.WriteHeader(http.StatusOK)
			}
		}
		mc := eb.ModelChanged{SourceId: model.Id}
		mcJson, _ := utils.SerializeJson(mc)
		dataChangedPub := eb.DataChangedEvent{Rb: rb, ModelName: "source", ParamsJson: mcJson, Op: eb.SAVE}
		dataChangedPub.Publish()

		config, _ := rb.ConfigRep.GetConfig()
		utils.CreateSourceDefaultDirectories(config, model.Id)

		ctx.JSON(http.StatusOK, model)
	})
	router.DELETE("/sources/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		if err := rb.SourceRep.RemoveById(id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//stops the stream after delete.
		ssrEvent := eb.StopStreamRequestEvent{Rb: rb, Id: id}
		err := ssrEvent.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}

		//also remove Object Detection Model
		if err := rb.OdRep.RemoveById(id); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//

		mc := eb.ModelChanged{SourceId: id}
		mcJson, _ := utils.SerializeJson(mc)
		dcEvent := eb.DataChangedEvent{Rb: rb, ModelName: "od", ParamsJson: mcJson, Op: eb.DELETE}
		dcEvent.Publish()

		dataChangedPub := eb.DataChangedEvent{Rb: rb, ModelName: "source", ParamsJson: mcJson, Op: eb.DELETE}
		dataChangedPub.Publish()

		ctx.JSON(http.StatusOK, gin.H{"id": id})
	})
	router.GET("/sourcestreamstatus", func(context *gin.Context) {
		modelList, err := rb.SourceRep.GetSourceStreamStatus(rb.StreamRep)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		context.JSON(http.StatusOK, modelList)
	})

	router.POST("/setsourceenabled", func(ctx *gin.Context) {
		var model models.SourceEnabledModel
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if len(model.SourceId) == 0 {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		source, err := rb.SourceRep.Get(model.SourceId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}
		if !model.Enabled {
			//stops the stream before disable.
			ssrEvent := eb.StopStreamRequestEvent{Rb: rb, Id: source.Id}
			ssrEvent.Publish()
			time.Sleep(time.Second * 5)
		}
		source.Enabled = model.Enabled
		rb.SourceRep.Save(source)
		ctx.JSON(http.StatusOK, source)
	})
}
