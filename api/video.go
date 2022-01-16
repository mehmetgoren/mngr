package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

func RegisterVideoEndpoints(router *gin.Engine) {
	router.GET("/videos/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		//id = "QLma6mWR3V8"
		files, _ := ioutil.ReadDir(utils.RelativePlaybackFolderPath + "/" + id)
		var list = make([]*models.VideoFile, 0)
		for _, file := range files {
			videoFile := models.VideoFile{}
			videoFile.SourceId = id
			videoFile.Name = file.Name()
			videoFile.Path = utils.RelativePlaybackFolderPath + "/" + file.Name()
			videoFile.Size = utils.Round(float64(file.Size()) * 0.000001)
			videoFile.CreatedAt = file.Name()
			videoFile.ModifiedAt = utils.FromDateToString(file.ModTime())
			list = append(list, &videoFile)
		}
		ctx.JSON(http.StatusOK, list)
	})

	router.GET("/recording/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		rep := reps.RecordingRepository{Connection: connSources}
		model, err := rep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, model)
	})
}
