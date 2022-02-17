package api

import (
	"bytes"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"mngr/models"
	"mngr/utils"
	"net/http"
	"os"
	"time"
)

func RegisterStreamEndpoints(router *gin.Engine) {
	router.GET("/stream", func(c *gin.Context) {
		modelList, err := utils.StreamRep.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		config, _ := utils.ConfigRep.GetConfig()
		for _, stream := range modelList {
			utils.SetHlsPath(config, stream)
		}

		c.JSON(http.StatusOK, modelList)
	})
	router.GET("/stream/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		stream, err := utils.StreamRep.Get(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		config, _ := utils.ConfigRep.GetConfig()
		utils.SetHlsPath(config, stream)
		ctx.JSON(http.StatusOK, stream)
	})
	router.POST("/ffmpegreader", func(context *gin.Context) {
		var readerModel models.FFmpegReaderModel
		if err := context.BindJSON(&readerModel); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//write base64 image to file
		b := readerModel.Img
		unbased, err := base64.StdEncoding.DecodeString(b)
		if err != nil {
			panic("Cannot decode b64")
		}

		r := bytes.NewReader(unbased)
		im, err := jpeg.Decode(r)
		if err != nil {
			panic("Bad jpeg")
		}

		f, err := os.OpenFile("/mnt/sdc1/pics/"+utils.FromDateToString(time.Now())+".jpg", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}

		o := jpeg.Options{Quality: 100}
		jpeg.Encode(f, im, &o)

		context.JSON(http.StatusOK, true)
	})
}
