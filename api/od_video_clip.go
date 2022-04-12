package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/view_models"
	"net/http"
	"os"
	"path"
)

func RegisterOdVideoClipEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/odvideoclips/:sourceid/:date", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()

		sourceId := ctx.Param("sourceid")
		date := ctx.Param("date")

		odhRep := reps.OdHandlerRepository{Config: config}
		jsonObjects := odhRep.GetJsonObjects(sourceId, date, true)
		list := view_models.Map(jsonObjects)

		ctx.JSON(http.StatusOK, list)
	})

	//potential security risk -> filename
	router.DELETE("/odvideoclips", func(ctx *gin.Context) {
		var vm view_models.OdVideoClipsViewModel
		if err := ctx.ShouldBindJSON(&vm); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(vm.DataFileNames) == 0 {
			ctx.JSON(http.StatusBadRequest, false)
		}

		config, _ := rb.ConfigRep.GetConfig()
		rootPath := config.General.RootFolderPath

		os.Remove(path.Join(rootPath, vm.VideoFileName))
		for _, ifn := range vm.ImageFileNames {
			os.Remove(path.Join(rootPath, ifn))
		}
		for _, dfi := range vm.DataFileNames {
			os.Remove(path.Join(rootPath, dfi))
		}

		ctx.JSON(http.StatusOK, true)
	})
}
