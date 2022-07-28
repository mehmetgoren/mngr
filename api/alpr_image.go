package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterAlprImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket, factory *cmn.Factory) {
	router.GET("/alprimagesfolders/:id/:date", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		date := ctx.Param("date")
		config, _ := rb.ConfigRep.GetConfig()
		odPath := utils.GetHourlyAlprImagesPathBySourceId(config, sourceId, date)
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("alprimages", func(ctx *gin.Context) {
		var model ImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items := make([]*ImageItem, 0)
		si := models.CreateDateSort(factory.GetCreatedDateFieldName())
		dtos, err := factory.CreateRepository().QueryAlprs(*data.GetParamsByHour(model.SourceId, model.RootPath, si))
		if err != nil {
			ctx.JSON(http.StatusOK, items)
			return
		}

		for _, dto := range dtos {
			item := &ImageItem{Id: dto.Id, ImagePath: dto.ImageFileName, CreatedAt: dto.CreatedAt}
			items = append(items, item)
		}

		ctx.JSON(http.StatusOK, items)
	})
}
