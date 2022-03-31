package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
)

func RegisterStaticResources(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	streamFolderPath, _ := utils.GetStreamFolderPath()
	router.Static("/livestream", streamFolderPath)
	recordFolderPath, _ := utils.GetRecordFolderPath()
	router.Static("/playback", recordFolderPath)
	detectedFolderName := utils.GetDetectedFolderName()
	router.Static("/"+detectedFolderName, "./static/"+detectedFolderName)
	router.Static("livestreamexample", "./static/live/example.mp4")
}
