package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterFrImagesEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/frimagesfolders/:id", func(ctx *gin.Context) {
		sourceId := ctx.Param("id")
		config, _ := rb.ConfigRep.GetConfig()
		odPath := utils.GetFrImagesPathBySourceId(config, sourceId)
		items, _ := newTree(odPath, true)
		ctx.JSON(http.StatusOK, items)
	})
	// it has potential security risk
	router.POST("frimages", func(ctx *gin.Context) {
		var model DetectedImagesParams
		if err := ctx.ShouldBindJSON(&model); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		config, _ := rb.ConfigRep.GetConfig()
		frhRep := reps.FrHandlerRepository{Config: config}
		jsonObjects := frhRep.GetJsonObjects(model.SourceId, model.RootPath, true)
		items := make([]*ImageItem, 0)
		for _, jsonObject := range jsonObjects {
			fr := jsonObject.FaceRecognition
			item := &ImageItem{Id: fr.Id, ImagePath: fr.ImageFileName}
			items = append(items, item)
		}
		ctx.JSON(http.StatusOK, items)
	})
}
