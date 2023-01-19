package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/utils"
	"strings"
)

var whiteList = make([]string, 0)

func initWhiteList() {
	whiteList = append(whiteList, "/blank.mp4")
}
func GetWhiteList() []string {
	return whiteList
}
func RegisterStaticResources(router *gin.Engine, rb *reps.RepoBucket) {
	config, _ := rb.ConfigRep.GetConfig()

	utils.CreateRequiredDirectories(config)

	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	router.StaticFile("/blank.mp4", "./static/playback/blank.mp4")

	initWhiteList()
	for _, dirPath := range config.General.DirPaths {
		relativePath := dirPath
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}
		router.Static(relativePath, dirPath)
		whiteList = append(whiteList, relativePath)
	}

	sources, err := rb.SourceRep.GetAll()
	if err == nil && sources != nil {
		for _, source := range sources {
			utils.CreateSourceDefaultDirectories(config, source.Id)
		}
	}
}
