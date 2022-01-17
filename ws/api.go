package ws

import (
	"github.com/gin-gonic/gin"
	"mngr/eb"
	"mngr/models"
	"net/http"
)

func RegisterApiEndpoints(router *gin.Engine) {
	router.POST("/startstreaming", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		eventPub := eb.StreamingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	router.POST("/startrecording", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		eventPub := eb.StartRecordingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	router.POST("/stoprecording", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		eventPub := eb.StopRecordingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
}

func RegisterWsEndpoints(router *gin.Engine) {
	//websockets
	router.StaticFile("/home", "./static/live/home.html")
	hub := NewHub()
	go hub.Run()
	router.GET("/wschat", func(ctx *gin.Context) {
		HandlerChat(hub, ctx.Writer, ctx.Request)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstreaming", func(ctx *gin.Context) {
		clientStreaming := CreateClient(hub, ctx.Writer, ctx.Request)
		eb.SubscribeStreamingEvents(clientStreaming)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstartrecording", func(ctx *gin.Context) {
		clientsStartRecording := CreateClient(hub, ctx.Writer, ctx.Request)
		eb.SubscribeStartRecordingEvents(clientsStartRecording)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstoprecording", func(ctx *gin.Context) {
		clientsStopRecording := CreateClient(hub, ctx.Writer, ctx.Request)
		eb.SubscribeStopRecordingEvents(clientsStopRecording)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	//end websockets
}
