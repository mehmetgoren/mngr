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
	"mngr/dsk_archv"
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
	err := factory.Init()
	if err != nil {
		log.Fatalln("factory init error: ", err.Error())
	}
	return factory
}

func main() {
	defer utils.HandlePanic()

	initWhiteList()
	rb.Init()
	global := ReadEnvVariables(rb)
	holders.Init()
	WhoAreYou(rb)
	FetchRtspTemplates(rb)

	factory := createFactory()
	defer func(factory *cmn.Factory) {
		err := factory.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(factory)
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
	dskChckr := dsk_archv.InitDiskUsageChecker(factory, rb)
	defer func() {
		dskChckr.StopScheduler()
	}()

	checkDefaultUser(rb)

	router := gin.Default()
	f, _ := os.Create("access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router.Use(gin.Logger())
	router.Use(authMiddleware)
	if global.ReadOnlyMode {
		router.Use(readonlyMiddleware)
	}
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
	api.RegisterOthersEndpoints(router, rb, global)
	api.RegisterCloudEndpoints(router, rb)
	api.RegisterAiDataEndpoints(router, factory)

	ws.RegisterApiEndpoints(router, rb)
	ws.RegisterWsEndpoints(router, holders)

	config, _ := rb.ConfigRep.GetConfig()
	log.Println("Root path is", config.General.RootFolderPath)
	port := strconv.Itoa(utils.ParsePort())
	log.Println("web server port is " + port)
	err = router.Run(":" + port)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

var anonymousEndPoints = map[string]int{"/login": 1, "/registeruser": 1, "/isReadOnlyMode": 1}

func authMiddleware(ctx *gin.Context) {
	req := ctx.Request
	if req.Method == "OPTIONS" {
		return
	}
	uri := req.RequestURI
	if _, ok := anonymousEndPoints[uri]; ok {
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
				err2 := ctx.AbortWithError(http.StatusUnauthorized, err)
				if err2 != nil {
					log.Println(err2.Error())
				}
				log.Println("an unauthorized request has been detected: ", req)
			}
		}
	}
}

var mutableEndPoints = map[string]int{"DELETE/aiclips": 1, "DELETE/deleteaidata": 1,
	"POST/telegram": 1, "DELETE/telegramuser": 1, "POST/gdrive": 1, "POST/resetgdrivetokenandurl": 1,
	"POST/config": 1, "GET/restoreconfig": 1,
	"DELETE/frtrainpersonimage": 1, "POST/frtrainpersonimage": 1, "POST/frtrainpersonrename": 1, "POST/frtrainpersonnew": 1, "DELETE/frtrainpersondelete": 1,
	"POST/ods":            1,
	"DELETE/records":      1,
	"POST/restartservice": 1, "POST/startservice": 1, "POST/stopservice": 1, "POST/restartaftercloudchanges": 1, "POST/restartallservices": 1,
	"POST/sources": 1, "DELETE/sources": 1, "POST/setsourceenabled": 1,
	"DELETE/users":     1,
	"POST/startstream": 1, "POST/stopstream": 1, "POST/restartstream": 1, "POST/videomerge": 1, "POST/frtrain": 1,
}

func readonlyMiddleware(ctx *gin.Context) {
	req := ctx.Request
	if req.Method == "OPTIONS" {
		return
	}
	var params = ctx.Params
	key := req.Method + req.RequestURI
	for _, param := range params {
		key = strings.Replace(key, "/"+param.Value, "", -1)
	}
	if _, ok := mutableEndPoints[key]; ok {
		err := errors.New("unauthorized 401")
		err2 := ctx.AbortWithError(http.StatusUnauthorized, err)
		if err2 != nil {
			log.Println(err2.Error())
		}
	}
}
