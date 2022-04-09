package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
)

func RegisterStaticResources(router *gin.Engine, rb *reps.RepoBucket) {
	config, _ := rb.ConfigRep.GetConfig()

	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	streamFolderPath, _ := utils.GetStreamFolderPath(config)
	router.Static("/livestream", streamFolderPath)
	recordFolderPath, _ := utils.GetRecordFolderPath(config)
	router.Static("/playback", recordFolderPath)
	router.Static("/"+utils.GetDetectedFolderName(), utils.GetDetectedFolderPath())
	utils.GetOdFolder(config)
	router.Static("livestreamexample", "./static/live/example.mp4")

	utils.CleanDetectedFolder()
}
