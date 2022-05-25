package api

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

func RegisterFrTrainingEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/frtrainpersons", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()
		trainPath := utils.GetFrTrainPath(config)
		directories, err := ioutil.ReadDir(trainPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		ret := make([]*models.FrTrainViewModel, 0)
		for _, dir := range directories {
			if !dir.IsDir() {
				continue
			}
			item := &models.FrTrainViewModel{Name: dir.Name()}
			subRoot := path.Join(trainPath, dir.Name())
			subDirectories, err := ioutil.ReadDir(subRoot)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			item.ImagePaths = make([]string, 0)
			for _, subDir := range subDirectories {
				item.ImagePaths = append(item.ImagePaths, path.Join("fr", "ml", "train", item.Name, subDir.Name()))
			}
			ret = append(ret, item)
		}
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].Name < ret[j].Name
		})
		ctx.JSON(http.StatusOK, ret)
	})

	// it has potential security risk
	router.GET("frtrainpersonimages/:person", func(ctx *gin.Context) {
		person := ctx.Param("person")
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		rootPath := utils.GetFrTrainPathByPerson(config, person)
		directories, err := ioutil.ReadDir(rootPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		rootPath = strings.Replace(rootPath, config.General.RootFolderPath+"/", "", -1)
		items := make([]*ImageItem, 0)
		for _, dir := range directories {
			name := dir.Name()
			item := &ImageItem{Id: person + "_" + name, ImagePath: path.Join(rootPath, name),
				CreatedAt: utils.TimeToString(dir.ModTime(), false)}
			items = append(items, item)
		}
		sort.Slice(items, func(i, j int) bool {
			t1 := utils.StringToTime(items[i].CreatedAt)
			t2 := utils.StringToTime(items[j].CreatedAt)
			return t1.After(t2)
		})
		ctx.JSON(http.StatusOK, items)
	})

	router.DELETE("frtrainpersonimage/:imgPath", func(ctx *gin.Context) {
		imgPath := ctx.Param("imgPath")
		barr, err := base64.StdEncoding.DecodeString(imgPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		s := string(barr)
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		p := u.Path
		fullPath := path.Join(config.General.RootFolderPath, p)
		fmt.Println(u)
		err = os.Remove(fullPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, true)
	})
}
