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
	"strconv"
	"strings"
)

func RegisterRecordEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/recordhours/:id/:datestr", func(ctx *gin.Context) {
		id := ctx.Param("id")
		dateStr := ctx.Param("datestr")
		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath := utils.GetHourlyRecordPathBySourceId(config, id, dateStr)
		files, _ := ioutil.ReadDir(recordFolderPath)
		var list = make([]string, 0)
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			list = append(list, file.Name())
		}
		ctx.JSON(http.StatusOK, list)
	})
	router.GET("/records/:id/:datestr/:hour", func(ctx *gin.Context) {
		id := ctx.Param("id")
		dateStr := ctx.Param("datestr")
		hour := ctx.Param("hour")
		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath := path.Join(utils.GetHourlyRecordPathBySourceId(config, id, dateStr), hour)
		files, _ := ioutil.ReadDir(recordFolderPath)
		date := utils.StringToTime(dateStr)
		var list = make([]*models.VideoFile, 0)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			videoFile := models.VideoFile{}
			videoFile.SourceId = id
			videoFile.Name = file.Name()
			videoFile.Year = strconv.Itoa(date.Year())
			videoFile.Month = utils.FixZero(int(date.Month()))
			videoFile.Day = utils.FixZero(date.Day())
			intHour, _ := strconv.Atoi(hour)
			videoFile.Hour = utils.FixZero(intHour)
			utils.SetVideoFilePath(&videoFile)
			videoFile.Size = utils.Round(float64(file.Size()) * 0.000001)
			videoFile.CreatedAt = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			videoFile.ModifiedAt = utils.TimeToString(file.ModTime(), true)
			list = append(list, &videoFile)
		}
		ctx.JSON(http.StatusOK, list)
	})

	router.DELETE("/records", func(ctx *gin.Context) {
		var vf models.VideoFile
		if err := ctx.ShouldBindJSON(&vf); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath := utils.GetRecordPath(config)
		err := os.Remove(utils.GetVideoFileAbsolutePath(&vf, recordFolderPath))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		} else {
			ctx.JSON(http.StatusOK, true)
		}
	})
}
