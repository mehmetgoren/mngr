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
	router.POST("/startstreaming", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.StartStreamingRequestEvent{SourceModel: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/stopstreaming", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.StopStreamingRequestEvent{Id: source.Id}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/editor", func(ctx *gin.Context) {
		var event eb.EditorRequestEvent
		ctx.BindJSON(&event)
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	// restart does not need to be subscribed to since it is only called by the restart_streaming_request which is just a proxy.
	router.POST("/restartstreaming", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.RestartStreamingRequestEvent{SourceModel: source}
		err := eventPub.Publish()
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
	router.GET("/wsstartstreaming", func(ctx *gin.Context) {
		clientStreaming := CreateClient(hub, ctx.Writer, ctx.Request)
		streamingEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "start_streaming_response"}
		streamingEventSub := eb.StartStreamingResponseEvent{Pusher: clientStreaming}
		go streamingEventBusSub.Subscribe(&streamingEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstreaming", func(ctx *gin.Context) {
		clientStreaming := CreateClient(hub, ctx.Writer, ctx.Request)
		streamingEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "stop_streaming_response"}
		streamingEventSub := eb.StopStreamingResponseEvent{Pusher: clientStreaming}
		go streamingEventBusSub.Subscribe(&streamingEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		clientEditor := CreateClient(hub, ctx.Writer, ctx.Request)
		editorEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "editor_response"}
		editorEvent := eb.EditorResponseEvent{Pusher: clientEditor}
		go editorEventBus.Subscribe(&editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	// End Subscribe
}
