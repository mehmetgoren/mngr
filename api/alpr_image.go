package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterAlprImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/alprimagesfolders/:id", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		config, _ := rb.ConfigRep.GetConfig()
		odPath := utils.GetAlprImagesPathBySourceId(config, sourceId)
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
		config, _ := rb.ConfigRep.GetConfig()
		arhRep := reps.AlprHandlerRepository{Config: config}
		jsonObjects := arhRep.GetJsonObjects(model.SourceId, model.RootPath, true)
		items := make([]*ImageItem, 0)
		for _, jsonObject := range jsonObjects {
			ar := jsonObject.AlprResults
			item := &ImageItem{Id: ar.Id, ImagePath: ar.ImageFileName}
			items = append(items, item)
		}
		ctx.JSON(http.StatusOK, items)
	})
}
