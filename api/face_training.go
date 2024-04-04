package api

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"io/ioutil"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
)

func RegisterFaceTrainingEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/facetrainpersons", func(ctx *gin.Context) {
		config, _ := rb.ConfigRep.GetConfig()
		trainPath := utils.GetFaceTrainPath(config)
		directories, err := ioutil.ReadDir(trainPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		ret := make([]*models.FaceTrainViewModel, 0)
		for _, dir := range directories {
			if !dir.IsDir() {
				continue
			}
			item := &models.FaceTrainViewModel{Name: dir.Name()}
			personDir := path.Join(trainPath, dir.Name())
			personImages, err := ioutil.ReadDir(personDir)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			item.ImagePaths = make([]string, 0)
			for _, personImage := range personImages {
				item.ImagePaths = append(item.ImagePaths, path.Join(personDir, personImage.Name()))
			}
			ret = append(ret, item)
		}
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].Name < ret[j].Name
		})
		ctx.JSON(http.StatusOK, ret)
	})

	// it has potential security risk
	router.GET("facetrainpersonimages/:person", func(ctx *gin.Context) {
		person := ctx.Param("person")
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		rootPath := utils.GetFaceTrainPathByPerson(config, person)
		directories, err := ioutil.ReadDir(rootPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

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

	router.DELETE("facetrainpersonimage/:imgPath", func(ctx *gin.Context) {
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

		err = os.Remove(u.Path)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.POST("facetrainpersonimage", func(ctx *gin.Context) {
		viewModel := models.FaceTrainScreenshotViewModel{}
		err := ctx.BindJSON(&viewModel)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		for _, base64Image := range viewModel.Base64Images {
			b, err := base64.StdEncoding.DecodeString(base64Image)
			if err != nil || len(b) == 0 {
				continue
			}
			img, err := jpeg.Decode(bytes.NewReader(b))
			if err != nil {
				continue
			}
			filename := path.Join(utils.GetFaceTrainPathByPerson(config, viewModel.Name), utils.NewId()+".jpg")
			f, err := os.Create(filename)
			if err != nil {
				continue
			}
			err = jpeg.Encode(f, img, nil)
			f.Close()
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.POST("facetrainpersonrename", func(ctx *gin.Context) {
		viewModel := models.FaceTrainRename{}
		err := ctx.BindJSON(&viewModel)
		if err != nil || len(viewModel.NewName) == 0 || len(viewModel.OriginalName) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		originalDirName := utils.GetFaceTrainPathByPerson(config, viewModel.OriginalName)
		if _, err := os.Stat(originalDirName); err != nil {
			if os.IsNotExist(err) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
		check := utils.IsDirNameValid(viewModel.NewName)
		if !check {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		newDirName := utils.GetFaceTrainPathByPerson(config, viewModel.NewName)
		err = os.Rename(originalDirName, newDirName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.POST("facetrainpersonnew", func(ctx *gin.Context) {
		viewModel := models.FaceTrainName{}
		err := ctx.BindJSON(&viewModel)
		if err != nil || len(viewModel.Name) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		check := utils.IsDirNameValid(viewModel.Name)
		if !check {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		dirName := utils.GetFaceTrainPathByPerson(config, viewModel.Name)
		err = os.Mkdir(dirName, 0777)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.DELETE("facetrainpersondelete", func(ctx *gin.Context) {
		viewModel := models.FaceTrainName{}
		if err := ctx.ShouldBindJSON(&viewModel); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		check := utils.IsDirNameValid(viewModel.Name)
		if !check {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		}
		config, err := rb.ConfigRep.GetConfig()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		dirName := utils.GetFaceTrainPathByPerson(config, viewModel.Name)
		err = os.RemoveAll(dirName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		ctx.JSON(http.StatusOK, true)
	})
}
