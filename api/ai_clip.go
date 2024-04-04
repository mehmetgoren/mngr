package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"mngr/view_models"
	"net/http"
	"os"
)

var emptyAiCLipViewModels = make([]*view_models.AiClipViewModel, 0)

func RegisterAiVideoClipEndpoints(router *gin.Engine, factory *cmn.Factory) {
	router.POST("/aiclips", func(ctx *gin.Context) {
		var params view_models.AiClipQueryViewModel
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
			return
		}

		si := models.CreateDateSort(factory.GetCreatedDateFieldName())
		aiDtos := make([]*data.AiDto, 0)

		entities, err := factory.CreateRepository().QueryAis(*data.GetParamsByHour(params.SourceId, params.Module, params.Date, si))
		if err != nil {
			ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
			return
		}
		for _, entity := range entities {
			aiDtos = append(aiDtos, entity)
		}

		list := view_models.Map(params.SourceId, aiDtos)
		ctx.JSON(http.StatusOK, list)
	})

	//potential security risk -> filename
	router.DELETE("/aiclips", func(ctx *gin.Context) {
		var vm view_models.AiClipViewModel
		if err := ctx.ShouldBindJSON(&vm); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(vm.SourceId) == 0 || len(vm.Ids) == 0 {
			ctx.JSON(http.StatusBadRequest, false)
			return
		}

		os.Remove(vm.VideoFileName)
		for _, ifn := range vm.ImageFileNames {
			os.Remove(ifn)
		}
		rep := factory.CreateRepository()
		for _, id := range vm.Ids {
			rep.RemoveAi(id)
		}

		ctx.JSON(http.StatusOK, true)
	})
}
