package ws

import (
	"github.com/gin-gonic/gin"
	"log"
	"mngr/eb"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
)

// RegisterApiEndpoints Publish Start
func RegisterApiEndpoints(router *gin.Engine, rb *reps.RepoBucket) {
	router.POST("/startstream", func(ctx *gin.Context) {
		var source models.SourceModel
		ctx.BindJSON(&source)
		config, _ := rb.ConfigRep.GetConfig()
		utils.CreateSourceDefaultDirectories(config, source.Id)
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
		var event eb.VfmRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	router.POST("/frtrain", func(ctx *gin.Context) {
		var event eb.FaceTrainRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	router.POST("/probe", func(ctx *gin.Context) {
		var event eb.ProbeRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})
}

// Publish End

// RegisterWsEndpoints Subscribe Start
func RegisterWsEndpoints(router *gin.Engine, holders *Holders) {
	hub := NewHub()
	go hub.Run()

	router.GET("/wsstartstream", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, StartStreamEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstream", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, StopStreamEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		requester := utils.GetQsValue(ctx, "requester")
		holders.RegisterEndPoint(hub, ctx, EditorEvent, requester)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsffmpegreader", func(ctx *gin.Context) {
		id := utils.GetQsValue(ctx, "id")
		if len(id) == 0 {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			log.Println("wsffmpegreader invalid source Id value")
			return
		}
		holders.RegisterEndPoint(hub, ctx, FFmpegReaderEvent, id)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsonvif", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, OnvifEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsvideomerge", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, VideoMergeEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsfrtrain", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, FrTrainEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsprobe", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, ProbeEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsnotifier", func(ctx *gin.Context) {
		holders.RegisterEndPoint(hub, ctx, NotifierEvent, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsuserlogout", func(ctx *gin.Context) {
		qs := ctx.Request.URL.Query()
		if val, ok := qs["token"]; ok {
			if len(val) > 0 {
				userToken := val[0]
				client := CreateClient(hub, ctx.Writer, ctx.Request)
				holders.UserLogin(userToken, client)
				ctx.Writer.WriteHeader(http.StatusOK)
				return
			}
		}
		ctx.Writer.WriteHeader(http.StatusNotFound)
	})
	// End Subscribe
}
