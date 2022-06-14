package api

import (
	"github.com/gin-gonic/gin"
	"mngr/reps"
	"mngr/server_stats"
	"net/http"
)

func RegisterServerStatsEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.GET("/serverstats", func(ctx *gin.Context) {

		stats := &server_stats.ServerStats{}
		err := stats.InitCpuInfos()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		err = stats.InitMemInfos()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		config, err := rb.ConfigRep.RestoreConfig()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		err = stats.InitDiskInfos(config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		err = stats.InitNetworkInfos()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, stats)
	})
}
