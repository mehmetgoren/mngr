package ws

import (
	"github.com/gin-gonic/gin"
	"log"
	"mngr/eb"
	"mngr/models"
	"mngr/reps"
	"net/http"
)

// RegisterApiEndpoints Publish Start
func RegisterApiEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.POST("/startstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.StartStreamRequestEvent{Rb: rb, SourceModel: source}
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
		eventPub := eb.StopStreamRequestEvent{Rb: rb, Id: source.Id}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/editor", func(ctx *gin.Context) {
		var event eb.EditorRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	// restart does not need to be subscribed to since it is only called by the restart_stream_request which is just a proxy.
	router.POST("/restartstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		eventPub := eb.RestartStreamRequestEvent{Rb: rb, SourceModel: source}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/onvif", func(ctx *gin.Context) {
		var e models.OnvifEvent
		ctx.BindJSON(&e)
		eventPub := eb.OnvifRequestEvent{Rb: rb, OnvifEvent: e}
		err := eventPub.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
	router.POST("/videomerge", func(ctx *gin.Context) {
		var event eb.VideMergeRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
}

// Publish End

type FFmpegReaderHolder struct {
	EventBus *eb.EventBus
	Client   *Client
	Event    *eb.FFmpegReaderResponseEvent
}

var ffmpegReaderDic = make(map[string]*FFmpegReaderHolder)

// RegisterWsEndpoints Subscribe Start
func RegisterWsEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.StaticFile("/home", "./static/live/home.html")
	hub := NewHub()
	go hub.Run()

	router.GET("/wsstartstream", func(ctx *gin.Context) {
		clientStream := CreateClient(hub, ctx.Writer, ctx.Request)
		streamEventBusSub := eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "start_stream_response"}
		streamEventSub := eb.StartStreamResponseEvent{Rb: rb, Pusher: clientStream}
		go streamEventBusSub.Subscribe(&streamEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstream", func(ctx *gin.Context) {
		clientStream := CreateClient(hub, ctx.Writer, ctx.Request)
		streamEventBusSub := eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "stop_stream_response"}
		streamEventSub := eb.StopStreamResponseEvent{Pusher: clientStream}
		go streamEventBusSub.Subscribe(&streamEventSub)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		clientEditor := CreateClient(hub, ctx.Writer, ctx.Request)
		editorEventBus := eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "editor_response"}
		editorEvent := eb.EditorResponseEvent{Pusher: clientEditor}
		go editorEventBus.Subscribe(&editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	//todo: add key to user id when login is available.
	router.GET("/wsffmpegreader", func(ctx *gin.Context) {
		qs := ctx.Request.URL.Query()
		if _, ok := qs["id"]; !ok {
			log.Println("wsffmpegreader invalid qs")
			return
		}
		id := qs["id"][0]
		if len(id) == 0 {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			log.Println("wsffmpegreader invalid qs")
			return
		}
		wsClient := CreateClient(hub, ctx.Writer, ctx.Request)

		if prev, ok := ffmpegReaderDic[id]; ok {
			err := prev.Client.Close()
			if err != nil {
				log.Println("Error while closing prev websockets connection for FFmPEG Reader. Err: ", err)
			}
			prev.Event.Pusher = wsClient
			log.Println("wsffmpegreader item has been already added,changing Ws Client for " + id)
			return
		}

		editorEventBus := &eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "ffrs" + id}
		editorEvent := &eb.FFmpegReaderResponseEvent{Pusher: wsClient}
		ffmpegReaderDic[id] = &FFmpegReaderHolder{EventBus: editorEventBus, Client: wsClient, Event: editorEvent}
		go editorEventBus.Subscribe(editorEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsonvif", func(ctx *gin.Context) {
		clientOnvif := CreateClient(hub, ctx.Writer, ctx.Request)
		onvifEventBus := eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "onvif_response"}
		onvifEvent := eb.OnvifResponseEvent{Pusher: clientOnvif}
		go onvifEventBus.Subscribe(&onvifEvent)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsvideomerge", func(ctx *gin.Context) {
		c := CreateClient(hub, ctx.Writer, ctx.Request)
		eventBus := eb.EventBus{PubSubConnection: rb.PubSubConnection, Channel: "vfm_response"}
		event := eb.VideMergeResponseEvent{Pusher: c}
		go eventBus.Subscribe(&event)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	// End Subscribe
}
