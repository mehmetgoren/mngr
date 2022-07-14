package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data/cmn"
	"mngr/reps"
	"mngr/utils"
	"mngr/view_models"
	"net/http"
	"os"
	"path"
)

func RegisterOdVideoClipEndpoints(router *gin.Engine, rb *reps.RepoBucket, factory *cmn.Factory) {
	router.GET("/odvideoclips/:sourceid/:date", func(ctx *gin.Context) {
		sourceId := ctx.Param("sourceid")
		date := ctx.Param("date")

		t := utils.StringToTime(date)
		ti := utils.TimeIndex{}
		ti.SetValuesFrom(&t)

		entities, err := factory.CreateRepository().GetOds(sourceId, &ti, true)
		if err != nil {
			ctx.JSON(http.StatusOK, make([]*view_models.OdVideoClipsViewModel, 0))
			return
		}

		list := view_models.Map(entities)
		ctx.JSON(http.StatusOK, list)
	})

	//potential security risk -> filename
	router.DELETE("/odvideoclips", func(ctx *gin.Context) {
		var vm view_models.OdVideoClipsViewModel
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
		os.Remove(path.Join(rootPath, vm.VideoFileName))
		for _, ifn := range vm.ImageFileNames {
			os.Remove(path.Join(rootPath, ifn))
		}
		rep := factory.CreateRepository()
		for _, id := range vm.Ids {
			rep.RemoveOd(id)
		}

		ctx.JSON(http.StatusOK, true)
	})
}
