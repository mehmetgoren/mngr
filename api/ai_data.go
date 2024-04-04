package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"net/http"
)

var emptyList = make([]*data.AiDto, 0)

func query(factory *cmn.Factory, params *data.QueryParams) []*data.AiDto {
	ais, err := factory.CreateRepository().QueryAis(*params)
	if err != nil {
		ais = make([]*data.AiDto, 0)
	}
	return ais
}

func mapFromAdvanced(p *models.QueryAiDataAdvancedParams) *data.QueryParams {
	params := &data.QueryParams{}
	params.SourceId = p.SourceId
	params.Module = p.Module
	params.SetupTimes(p.StartDateTimeStr, p.EndDateTimeStr)
	params.Label = p.Label
	params.NoPreparingVideoFile = p.NoPreparingVideoFile
	params.Sort = p.Sort
	params.Paging = p.Paging

	return params
}

func RegisterAiDataEndpoints(router *gin.Engine, factory *cmn.Factory) {
	router.POST("queryaidataadvanced", func(ctx *gin.Context) {
		var p models.QueryAiDataAdvancedParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}
		params := mapFromAdvanced(&p)
		ctx.JSON(http.StatusOK, query(factory, params))
	})

	router.POST("queryaidatacount", func(ctx *gin.Context) {
		var p models.QueryAiDataAdvancedParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}

		params := mapFromAdvanced(&p)
		count, _ := factory.CreateRepository().CountAis(*params)
		ctx.JSON(http.StatusOK, count)
	})

	router.DELETE("deleteaidata", func(ctx *gin.Context) {
		var p models.AiDataDeleteOptions
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, false)
			return
		}
		options := &data.DeleteOptions{}
		options.Id = p.Id
		options.DeleteVideo = p.DeleteVideo
		options.DeleteImage = p.DeleteImage

		err := factory.CreateRepository().DeleteAis(options)
		if err != nil {
			ctx.JSON(http.StatusOK, false)
			return
		}

		ctx.JSON(http.StatusOK, true)
	})
}
