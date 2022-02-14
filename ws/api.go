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
	router.POST("/startstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.StartStreamRequestEvent{SourceModel: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/stopstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		if len(source.Id) == 0 {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		eventPub := eb.StopStreamRequestEvent{Id: source.Id}
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
	// restart does not need to be subscribed to since it is only called by the restart_stream_request which is just a proxy.
	router.POST("/restartstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.RestartStreamRequestEvent{SourceModel: source}
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
	router.GET("/wsstartstream", func(ctx *gin.Context) {
		clientStream := CreateClient(hub, ctx.Writer, ctx.Request)
		streamEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "start_stream_response"}
		streamEventSub := eb.StartStreamResponseEvent{Pusher: clientStream}
		go streamEventBusSub.Subscribe(&streamEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstream", func(ctx *gin.Context) {
		clientStream := CreateClient(hub, ctx.Writer, ctx.Request)
		streamEventBusSub := eb.EventBus{Connection: utils.ConnPubSub, Channel: "stop_stream_response"}
		streamEventSub := eb.StopStreamResponseEvent{Pusher: clientStream}
		go streamEventBusSub.Subscribe(&streamEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		clientEditor := CreateClient(hub, ctx.Writer, ctx.Request)
		editorEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "editor_response"}
		editorEvent := eb.EditorResponseEvent{Pusher: clientEditor}
		go editorEventBus.Subscribe(&editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsffmpegreader", func(ctx *gin.Context) {
		clientEditor := CreateClient(hub, ctx.Writer, ctx.Request)
		editorEventBus := eb.EventBus{Connection: utils.ConnPubSub, Channel: "read_service"}
		editorEvent := eb.FFmpegReaderResponseEvent{Pusher: clientEditor}
		go editorEventBus.Subscribe(&editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	// End Subscribe
}
