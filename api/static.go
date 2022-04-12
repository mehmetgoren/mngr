package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
)

func RegisterStaticResources(router *gin.Engine, rb *reps.RepoBucket) {
	config, _ := rb.ConfigRep.GetConfig()

	utils.CreateRequiredDirectories(config)

	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	streamFolderPath := utils.GetStreamPath(config)
	router.Static("/livestream", streamFolderPath)

	recordFolderPath := utils.GetRecordPath(config)
	router.Static("/playback", recordFolderPath)

	// od is not used here since od images are loaded from a temp file since the gallery performs better on non hierarchical filename instead of base64 strings
	router.Static("/od", utils.GetOdPath(config))
}
