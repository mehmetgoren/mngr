package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mngr/models"
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
			videoFile.Size = file.Size()
			videoFile.CreatedAt = file.Name()
			videoFile.ModifiedAt = utils.FromDateToString(file.ModTime())
			list = append(list, &videoFile)
		}
		ctx.JSON(http.StatusOK, list)
	})
}
