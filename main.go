package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"mngr/eb"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"mngr/ws"
	"net/http"
	"os"
	"strconv"
)

const (
	MAIN     = 0
	SERVICE  = 1
	SOURCES  = 2
	EVENTBUS = 3
)

func createRedisConnection(db int) *redis.Client {
	host := os.Getenv("REDIS_HOST")
	fmt.Println("Redis host: ", host)
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	portStr := os.Getenv("REDIS_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println("An error occurred while converting Redis port value:" + err.Error())
		port = 6379
	}

	return redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: "", // no password set
		DB:       db, // use default DB
	})
}

func main() {
	connSources := createRedisConnection(SOURCES)
	sourceRep := reps.SourceRepository{Connection: connSources}

	connConfig := createRedisConnection(MAIN)
	configRep := reps.ConfigRepository{Connection: connConfig}

	// subscribe to redis channel
	connPubSub := createRedisConnection(EVENTBUS)
	eventBusSub := eb.EventBus{Connection: connPubSub, Channel: "streaming_response"}
	eventSub := eb.StreamingEvent{}
	go eventBusSub.Subscribe(&eventSub)

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

	router.StaticFile("/favicon.ico", "./static/icons/favicon.ico")
	router.Static("/livestream", utils.RelativeLiveFolderPath)
	router.Static("livestreamexample", "./static/live/example.mp4")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/sources", func(ctx *gin.Context) {
		sources, _ := sourceRep.GetAllSources()
		ctx.JSON(http.StatusOK, sources)
	})

	router.GET("/mlconfig", func(ctx *gin.Context) {
		config, _ := configRep.GetMlConfig()
		ctx.JSON(http.StatusOK, config)
	})

	router.POST("/mlconfig", func(ctx *gin.Context) {
		var config models.MlConfig
		ctx.BindJSON(&config)
		configRep.SaveMlConfig(&config)
		ctx.JSON(http.StatusOK, config)
	})

	router.GET("/restoremlconfig", func(ctx *gin.Context) {
		config, _ := configRep.RestoreMlConfig()
		ctx.JSON(http.StatusOK, config)
	})

	router.POST("/startstreaming", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		folderPath, err := utils.CreateDirIfNotExist(utils.LiveFolderPath + "/" + source.Id)
		if err != nil {
			log.Println("An error occurred while creating folder: " + err.Error())
			return
		}

		folderPathFull, _ := utils.GetExecutablePath()
		eventPub := eb.StreamingEvent{Source: source, FolderPath: folderPathFull + "/" + folderPath + "/stream.m3u8"}
		eventBusPub := eb.EventBus{Connection: connPubSub, Channel: "streaming_request"}
		err = eventBusPub.Publish(&eventPub)
		if err != nil {
			log.Println("An error occurred while publishing event: " + err.Error())
			log.Println("panic you fuck")
			panic(err)
			return
		}

		ctx.Writer.WriteHeader(http.StatusOK)
	})

	//websockets
	router.StaticFile("/home", "./static/live/home.html")
	hub := ws.NewHub()
	go hub.Run()
	router.GET("/wschat", func(ctx *gin.Context) {
		ws.HandlerChat(hub, ctx.Writer, ctx.Request)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstreaming", func(ctx *gin.Context) {
		ws.HandlerStreaming(hub, ctx.Writer, ctx.Request)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	//end websockets

	router.Run(":2072")
}

func loggingMiddleware(ctx *gin.Context) {

}
