package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"mngr/reps"
	"mngr/view_models"
	"net/http"
	"os"
	"path"
	"strconv"
)

var emptyAiCLipViewModels = make([]*view_models.AiClipViewModel, 0)

func RegisterOdVideoClipEndpoints(router *gin.Engine, rb *reps.RepoBucket, factory *cmn.Factory) {
	router.POST("/aiclips", func(ctx *gin.Context) {
		var params view_models.AiClipQueryViewModel
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
			return
		}

		si := models.CreateDateSort(factory.GetCreatedDateFieldName())
		interfaces := make([]data.AiDto, 0)

		switch params.AiType {
		case models.Od:
			entities, err := factory.CreateRepository().QueryOds(*data.GetParamsByHour(params.SourceId, params.Date, si))
			if err != nil {
				ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
				return
			}
			for _, entity := range entities {
				interfaces = append(interfaces, entity)
			}
			break
		case models.Fr:
			entities, err := factory.CreateRepository().QueryFrs(*data.GetParamsByHour(params.SourceId, params.Date, si))
			if err != nil {
				ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
				return
			}
			for _, entity := range entities {
				interfaces = append(interfaces, entity)
			}
			break
		case models.Alpr:
			entities, err := factory.CreateRepository().QueryAlprs(*data.GetParamsByHour(params.SourceId, params.Date, si))
			if err != nil {
				ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
				return
			}
			for _, entity := range entities {
				interfaces = append(interfaces, entity)
			}
			break
		default:
			log.Println("an unsupported ai type has been found: " + strconv.Itoa(params.AiType))
			ctx.JSON(http.StatusOK, emptyAiCLipViewModels)
			return
		}

		list := view_models.Map(interfaces)
		ctx.JSON(http.StatusOK, list)
	})

	//potential security risk -> filename
	router.DELETE("/aiclips", func(ctx *gin.Context) {
		var vm view_models.AiClipViewModel
		if err := ctx.ShouldBindJSON(&vm); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(vm.Ids) == 0 {
			ctx.JSON(http.StatusBadRequest, false)
			return
		}

		config, _ := rb.ConfigRep.GetConfig()
		rootPath := config.General.RootFolderPath
		vm.RollbackVideoFileName()
		os.Remove(path.Join(rootPath, vm.VideoFileName))
		for _, ifn := range vm.ImageFileNames {
			os.Remove(path.Join(rootPath, ifn))
		}
		rep := factory.CreateRepository()
		for _, id := range vm.Ids {
			switch vm.AiType {
			case models.Od:
				rep.RemoveOd(id)
				break
			case models.Fr:
				rep.RemoveFr(id)
				break
			case models.Alpr:
				rep.RemoveAlpr(id)
				break
			default:
				log.Println("an unsupported ai type has been found: " + strconv.Itoa(vm.AiType))
				ctx.JSON(http.StatusOK, false)
				return
			}
		}

		ctx.JSON(http.StatusOK, true)
	})
}
