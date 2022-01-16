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

		eventPub := eb.RecordingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	//router.POST("/stopstreaming", func(ctx *gin.Context) { ...}
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
	router.GET("/wsrecording", func(ctx *gin.Context) {
		clientsRecording := CreateClient(hub, ctx.Writer, ctx.Request)
		eb.SubscribeRecordingEvents(clientsRecording)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	//end websockets
}
