package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
)

func RegisterStaticResources(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	streamingFolderPath, _ := utils.GetStreamingFolderPath()
	router.Static("/livestream", streamingFolderPath)
	recordingFolderPath, _ := utils.GetRecordingFolderPath()
	router.Static("/playback", recordingFolderPath)
	router.Static("livestreamexample", "./static/live/example.mp4")
}
