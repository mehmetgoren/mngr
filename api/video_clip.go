package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/fs"
	"io/ioutil"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RegisterVideoClipEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/videoclips/:sourceid/:date", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()

		sourceId := ctx.Param("sourceid")
		date := ctx.Param("date")
		t := utils.StringToDate(date)

		clips := make([]*models.VideoClipJsonObject, 0)
		recordFolderPath, _ := utils.GetRecordFolderPath(config)
		folderPath := path.Join(recordFolderPath, sourceId, "vcs", strconv.Itoa(t.Year()),
			strconv.Itoa(int(t.Month())), strconv.Itoa(t.Day()))

		filepath.Walk(folderPath, func(fileInfoFullPath string, file fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if file.IsDir() {
				return err
			}

			ext := path.Ext(file.Name())
			if ext != ".json" {
				return err
			}

			bytes, _ := ioutil.ReadFile(fileInfoFullPath)
			clip := &models.VideoClipJsonObject{}
			utils.DeserializeJsonB(bytes, clip)
			setupDateTimes(clip)
			clips = append(clips, clip)

			return nil
		})
		sort.Slice(clips, func(i, j int) bool {
			return clips[i].CreatedAtTime.After(clips[j].CreatedAtTime)
		})
		ctx.JSON(http.StatusOK, clips)
	})

	//potential security risk -> filename
	router.DELETE("/videoclips/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		fileNamBytes, err := base64.StdEncoding.DecodeString(id)
		fileName := string(fileNamBytes)

		config, _ := rb.ConfigRep.GetConfig()
		recordFolderPath, _ := utils.GetRecordFolderPath(config)
		err = os.Remove(path.Join(recordFolderPath, fileName))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		err = os.Remove(path.Join(recordFolderPath, strings.Replace(fileName, "mp4", "json", -1)))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, true)
	})
}

func setupDateTimes(v *models.VideoClipJsonObject) {
	fileName := utils.GetFileNameWithoutExtension(v.FileName)
	v.CreatedAtTime = utils.StringToTime(strings.Split(fileName, ".")[0], false)
	v.CreatedAt = utils.TimeToString(v.CreatedAtTime, false)
	v.LastModifiedTime = v.CreatedAtTime.Add(time.Duration(v.Duration * int(time.Second)))
	v.LastModified = utils.TimeToString(v.LastModifiedTime, false)
}
