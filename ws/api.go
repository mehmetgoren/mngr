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
	router.POST("/videomerge", func(ctx *gin.Context) {
		var event eb.VfmRequestEvent
		ctx.BindJSON(&event)
		event.Rb = rb
		err := event.Publish()
		if err == nil {
			ctx.Writer.WriteHeader(http.StatusOK)
		}
	})

	router.POST("/facetrain", func(ctx *gin.Context) {
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
func RegisterWsEndpoints(router *gin.Engine, hldrs *Holders) {
	hub := NewHub()
	go hub.Run()

	router.GET("/wsstartstream", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, StartStream, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsstopstream", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, StopStream, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wseditor", func(ctx *gin.Context) {
		requester := utils.GetQsValue(ctx, "requester")
		hldrs.RegisterEndPoint(hub, ctx, Editor, requester)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsffmpegreader", func(ctx *gin.Context) {
		id := utils.GetQsValue(ctx, "id")
		if len(id) == 0 {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			log.Println("wsffmpegreader invalid source Id value")
			return
		}
		hldrs.RegisterEndPoint(hub, ctx, FFmpegReader, id)
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsvideomerge", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, VideoMerge, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsfacetrain", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, FaceTrain, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsprobe", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, Probe, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsnotifier", func(ctx *gin.Context) {
		hldrs.RegisterEndPoint(hub, ctx, Notifier, "")
		ctx.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/wsuserlogout", func(ctx *gin.Context) {
		qs := ctx.Request.URL.Query()
		if val, ok := qs["token"]; ok {
			if len(val) > 0 {
				userToken := val[0]
				client := CreateClient(hub, ctx.Writer, ctx.Request)
				hldrs.UserLogin(userToken, client)
				ctx.Writer.WriteHeader(http.StatusOK)
				return
			}
		}
		ctx.Writer.WriteHeader(http.StatusNotFound)
	})
	// End Subscribe
}
