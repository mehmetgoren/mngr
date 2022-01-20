package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"mngr/api"
	"mngr/ws"
	"net/http"
	"os"
)

func main() {
	//utils.RemovePrevStreamFolders()
	router := gin.Default()
	f, _ := os.Create("access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router.Use(gin.Logger())
	router.Use(loggingMiddleware)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.RegisterStaticResources(router)
	api.RegisterSourceEndpoints(router)
	api.RegisterConfigEndpoints(router)
	api.RegisterVideoEndpoints(router)

	ws.RegisterApiEndpoints(router)
	ws.RegisterWsEndpoints(router)

	router.Run(":2072")
}

func loggingMiddleware(ctx *gin.Context) {

}
