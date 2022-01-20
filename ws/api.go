package ws

import (
	"github.com/gin-gonic/gin"
	"mngr/eb"
	"mngr/models"
	"mngr/utils"
	"net/http"
)

// Publish Start
func RegisterApiEndpoints(router *gin.Engine) {
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
	router.POST("/startstreaming", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		eventPub := eb.StartStreamingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/stopstreaming", func(ctx *gin.Context) {
		var source models.Source
		ctx.BindJSON(&source)
		eventPub := eb.StopStreamingEvent{Source: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/editor", func(ctx *gin.Context) {
		var event eb.EditorEvent
		ctx.BindJSON(&event)
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
}

// Publish End

// Subscribe Start
func RegisterWsEndpoints(router *gin.Engine) {
	router.StaticFile("/home", "./static/live/home.html")
	hub := NewHub()
	go hub.Run()
	router.GET("/wschat", func(ctx *gin.Context) {
		HandlerChat(hub, ctx.Writer, ctx.Request)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstartrecording", func(ctx *gin.Context) {
		clientsStartRecording := CreateClient(hub, ctx.Writer, ctx.Request)
		recordingEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "start_recording_response"}
		recordingEvent := eb.StartRecordingEvent{Pusher: clientsStartRecording}
		go recordingEventBus.Subscribe(&recordingEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstoprecording", func(ctx *gin.Context) {
		clientsStopRecording := CreateClient(hub, ctx.Writer, ctx.Request)
		recordingEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "stop_recording_response"}
		recordingEvent := eb.StopRecordingEvent{Pusher: clientsStopRecording}
		go recordingEventBus.Subscribe(&recordingEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstartstreaming", func(ctx *gin.Context) {
		clientStreaming := CreateClient(hub, ctx.Writer, ctx.Request)
		streamingEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "start_streaming_response"}
		streamingEventSub := eb.StartStreamingEvent{Pusher: clientStreaming}
		go streamingEventBusSub.Subscribe(&streamingEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstreaming", func(ctx *gin.Context) {
		clientStreaming := CreateClient(hub, ctx.Writer, ctx.Request)
		streamingEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "stop_streaming_response"}
		streamingEventSub := eb.StopStreamingEvent{Pusher: clientStreaming}
		go streamingEventBusSub.Subscribe(&streamingEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		clientEditor := CreateClient(hub, ctx.Writer, ctx.Request)
		editorEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "editor_response"}
		editorEvent := eb.EditorEvent{Pusher: clientEditor}
		go editorEventBus.Subscribe(&editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	// End Subscribe
}
