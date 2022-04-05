package api

import (
	"github.com/gin-gonic/gin"
	"mngr/eb"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"path"
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

		config, _ := rb.ConfigRep.GetConfig()
		// Create HLS stream folder.
		sf, _ := utils.GetStreamFolderPath(config)
		dic := path.Join(sf, model.Id)
		utils.CreateDicIfNotExist(dic)

		// Create record folder.
		rf, _ := utils.GetRecordFolderPath(config)
		dic = path.Join(rf, model.Id)
		utils.CreateDicIfNotExist(dic)
		//and also short video clips folder
		vcsDicPath := path.Join(dic, "vcs")
		utils.CreateDicIfNotExist(vcsDicPath)
		tempSvcDicPath := path.Join(vcsDicPath, "temp")
		utils.CreateDicIfNotExist(tempSvcDicPath)

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
		err = dcEvent.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
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
}
