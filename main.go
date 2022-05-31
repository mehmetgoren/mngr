package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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

	users, err := rb.UserRep.GetUsers()
	if users != nil {
		log.Println("user count: ", len(users))
	}

	router := gin.Default()
	f, _ := os.Create("access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router.Use(gin.Logger())
	router.Use(authMiddleware)
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
	api.RegisterOdImagesEndpoints(router, &rb)
	api.RegisterOdVideoClipEndpoints(router, &rb)
	api.RegisterFrImagesEndpoints(router, &rb)
	api.RegisterAlprImagesEndpoints(router, &rb)
	api.RegisterOnvifEndpoints(router, &rb)
	api.RegisterFrTrainingEndpoints(router, &rb)
	api.RegisterUserEndpoints(router, &rb)

	ws.RegisterApiEndpoints(router, &rb)
	ws.RegisterWsEndpoints(router, &rb)

	err = router.Run(":2072")
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func authMiddleware(ctx *gin.Context) {

}
