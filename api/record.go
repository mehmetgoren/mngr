package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func RegisterRecordEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/records/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath, _ := utils.GetRecordFolderPath(config)
		files, _ := ioutil.ReadDir(path.Join(recordFolderPath, id))
		var list = make([]*models.VideoFile, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			videoFile := models.VideoFile{}
			videoFile.SourceId = id
			videoFile.Name = file.Name()
			videoFile.Path = path.Join("/playback", id, file.Name()) //utils.RelativePlaybackFolderPath + "/" + file.Name()
			videoFile.Size = utils.Round(float64(file.Size()) * 0.000001)
			videoFile.CreatedAt = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			videoFile.ModifiedAt = utils.TimeToString(file.ModTime(), true)
			list = append(list, &videoFile)
		}
		ctx.JSON(http.StatusOK, list)
	})

	router.DELETE("/records/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		var fileNames []string
		ctx.BindJSON(&fileNames)
		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath, _ := utils.GetRecordFolderPath(config)
		for _, fileName := range fileNames {
			err := os.Remove(path.Join(recordFolderPath, id, fileName))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
		}
		ctx.JSON(http.StatusOK, nil)
	})
}
