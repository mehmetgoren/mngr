package main

import (
	"errors"
	"github.com/docker/docker/client"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mngr/api"
	"mngr/data/cmn"
	"mngr/dsk_usg"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"mngr/ws"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func checkDefaultUser(rb *reps.RepoBucket) {
	users, err := rb.UserRep.GetUsers()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if users != nil && len(users) > 0 {
		log.Println("user count: ", len(users))
	} else {
		_, err = rb.UserRep.Register(&models.RegisterUserViewModel{
			Username:   "admin",
			Password:   "admin",
			RePassword: "admin",
			Email:      "admin@feniks.com",
		})
		if err != nil {
			log.Println(err.Error())
			return
		}
		checkDefaultUser(rb)
	}
}

var rb = &reps.RepoBucket{}
var holders = &ws.Holders{Rb: rb}
var whiteList = make([]string, 0)

func initWhiteList() {
	whiteList = append(whiteList, "/livestream/")
	whiteList = append(whiteList, "/playback/")
	whiteList = append(whiteList, "/od/")
	whiteList = append(whiteList, "/fr/")
	whiteList = append(whiteList, "/alpr/")
	whiteList = append(whiteList, "/blank.mp4")
}

func createFactory() *cmn.Factory {
	config, _ := rb.ConfigRep.GetConfig()
	factory := &cmn.Factory{Config: config}
	factory.Init()
	return factory
}

func main() {
	defer utils.HandlePanic()

	initWhiteList()
	rb.Init()
	holders.Init()
	WhoAreYou(rb)

	factory := createFactory()
	defer factory.Close()
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println(err.Error())
	}
	defer func() {
		if dockerClient != nil {
			err = dockerClient.Close()
			log.Println(err.Error())
		}
	}()
	dskChckr := dsk_usg.InitDiskUsageChecker(factory, rb)
	defer func() {
		dskChckr.StopScheduler()
	}()

	checkDefaultUser(rb)

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

	api.RegisterStaticResources(router, rb)
	api.RegisterSourceEndpoints(router, rb)
	api.RegisterStreamEndpoints(router, rb)
	api.RegisterConfigEndpoints(router, rb)
	api.RegisterRecordEndpoints(router, rb)
	api.RegisterOdEndpoints(router, rb)
	api.RegisterOdImagesEndpoints(router, rb, factory)
	api.RegisterOdVideoClipEndpoints(router, rb, factory)
	api.RegisterOnvifEndpoints(router, rb)
	api.RegisterFrTrainingEndpoints(router, rb)
	api.RegisterUserEndpoints(router, holders)
	api.RegisterServiceEndpoints(router, rb, dockerClient)
	api.RegisterServerStatsEndpoints(router, rb)
	api.RegisterOthersEndpoints(router, rb)
	api.RegisterCloudEndpoints(router, rb)
	api.RegisterAiDataEndpoints(router, factory)

	ws.RegisterApiEndpoints(router, rb)
	ws.RegisterWsEndpoints(router, holders)

	port := strconv.Itoa(utils.ParsePort())
	log.Println("web server port is " + port)
	err = router.Run(":" + port)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func authMiddleware(ctx *gin.Context) {
	req := ctx.Request
	uri := req.RequestURI
	if req.Method == "OPTIONS" || uri == "/login" || uri == "/registeruser" {
		return
	}
	if strings.HasPrefix(uri, "/ws") { // if it is a websocket request
		qs := ctx.Request.URL.Query()
		if _, ok := qs["token"]; !ok {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			ctx.AbortWithStatus(http.StatusBadRequest)
			log.Println("websocket invalid query string parameters")
			return
		}
		token := qs["token"][0]
		if len(token) == 0 {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			ctx.AbortWithStatus(http.StatusBadRequest)
			log.Println("websocket invalid query string parameters")
			return
		}

		if _, ok := rb.IsUserAuthenticated(token); !ok {
			err := errors.New("unauthorized 401")
			holders.UserLogout(token, true)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
			log.Println("websocket token was not found")
		}
	} else {
		token := ctx.Request.Header.Get("user")
		if _, ok := rb.IsUserAuthenticated(token); !ok {
			isWhitelisted := false
			for _, wl := range whiteList {
				if strings.HasPrefix(uri, wl) {
					isWhitelisted = true
					break
				}
			}
			if !isWhitelisted {
				err := errors.New("unauthorized 401")
				holders.UserLogout(token, true)
				ctx.AbortWithError(http.StatusUnauthorized, err)
				log.Println("an unauthorized request has been detected: ", req)
			}
		}
	}
}
