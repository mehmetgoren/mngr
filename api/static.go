package api

import (
	"github.com/gin-gonic/gin"
	"mngr/utils"
)

func RegisterStaticResources(router *gin.Engine) {
	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	router.Static("/livestream", utils.RelativeLiveFolderPath)
	router.Static("/playback", utils.RelativePlaybackFolderPath)
	router.Static("livestreamexample", "./static/live/example.mp4")
}
