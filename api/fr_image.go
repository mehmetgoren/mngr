package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterFrImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket, factory *cmn.Factory) {
	router.GET("/frimagesfolders/:id/:date", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		date := ctx.Param("date")
		config, _ := rb.ConfigRep.GetConfig()
		odPath := utils.GetHourlyFrImagesPathBySourceId(config, sourceId, date)
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("frimages", func(ctx *gin.Context) {
		var model ImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		items := make([]*ImageItem, 0)
		dtos, err := factory.CreateRepository().GetFrs(data.GetParamsByHour(model.SourceId, model.RootPath))
		if err != nil {
			ctx.JSON(http.StatusOK, items)
			return
		}

		config, _ := rb.ConfigRep.GetConfig()
		for _, dto := range dtos {
			item := &ImageItem{Id: dto.Id, ImagePath: utils.SetRelativeImagePath(config, dto.ImageFileName), CreatedAt: dto.CreatedAt}
			items = append(items, item)
		}
		ctx.JSON(http.StatusOK, items)
	})
}
