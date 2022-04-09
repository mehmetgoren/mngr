package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"mngr/api"
	"mngr/reps"
	"mngr/utils"
	"mngr/ws"
	"net/http"
	"os"
)

func main() {
	defer utils.HandlePanic()

	rb := reps.RepoBucket{}
	rb.Init()

	WhoAreYou(&rb)

	router := gin.Default()
	f, _ := os.Create("access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router.Use(gin.Logger())
	router.Use(loggingMiddleware)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.RegisterStaticResources(router, &rb)
	api.RegisterSourceEndpoints(router, &rb)
	api.RegisterStreamEndpoints(router, &rb)
	api.RegisterConfigEndpoints(router, &rb)
	api.RegisterRecordEndpoints(router, &rb)
	api.RegisterOdEndpoints(router, &rb)
	api.RegisterDetectedEndpoints(router, &rb)
	api.RegisterVideoClipEndpoints(router, &rb)

	ws.RegisterApiEndpoints(router, &rb)
	ws.RegisterWsEndpoints(router, &rb)

	router.Run(":2072")
}

func loggingMiddleware(ctx *gin.Context) {

}
