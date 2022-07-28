package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"net/http"
)

var emptyList = make([]*data.AiDataDto, 0)

func query(aiType int, factory *cmn.Factory, params *data.QueryParams) []*data.AiDataDto {
	ret := make([]*data.AiDataDto, 0)
	switch aiType {
	case models.Od:
		ods, _ := factory.CreateRepository().QueryOds(*params)
		for _, od := range ods {
			ai := &data.AiDataDto{}
			ai.MapFromOd(od)
			ret = append(ret, ai)
		}
		break
	case models.Fr:
		frs, _ := factory.CreateRepository().QueryFrs(*params)
		for _, fr := range frs {
			ai := &data.AiDataDto{}
			ai.MapFromFr(fr)
			ret = append(ret, ai)
		}
		break
	case models.Alpr:
		alprs, _ := factory.CreateRepository().QueryAlprs(*params)
		for _, alpr := range alprs {
			ai := &data.AiDataDto{}
			ai.MapFromAlpr(alpr)
			ret = append(ret, ai)
		}
		break
	}
	return ret
}

func mapFromAdvanced(p *models.QueryAiDataAdvancedParams) *data.QueryParams {
	params := &data.QueryParams{}
	params.SourceId = p.SourceId
	params.SetupTimes(p.StartDateTimeStr, p.EndDateTimeStr)
	params.ClassName = p.PredClassName
	params.NoPreparingVideoFile = p.NoPreparingVideoFile
	params.Sort = p.Sort
	params.Paging = p.Paging

	return params
}

func RegisterAiDataEndpoints(router *gin.Engine, factory *cmn.Factory) {
	router.POST("queryaidata", func(ctx *gin.Context) {
		var p models.QueryAiDataParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}

		params := &data.QueryParams{}
		params.SourceId = p.SourceId
		params.ClassName = p.PredClassName
		params.Sort = models.CreateDateSort(factory.GetCreatedDateFieldName())
		params.NoPreparingVideoFile = p.NoPreparingVideoFile
		params.SetupHourlyTimes(p.DateTimeStr)

		ctx.JSON(http.StatusOK, query(p.AiType, factory, params))
	})

	router.POST("queryaidataadvanced", func(ctx *gin.Context) {
		var p models.QueryAiDataAdvancedParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}
		params := mapFromAdvanced(&p)
		ctx.JSON(http.StatusOK, query(p.AiType, factory, params))
	})

	router.POST("queryaidatacount", func(ctx *gin.Context) {
		var p models.QueryAiDataAdvancedParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}

		params := mapFromAdvanced(&p)
		switch p.AiType {
		case models.Od:
			count, _ := factory.CreateRepository().CountOds(*params)
			ctx.JSON(http.StatusOK, count)
			break
		case models.Fr:
			count, _ := factory.CreateRepository().CountFrs(*params)
			ctx.JSON(http.StatusOK, count)
			break
		case models.Alpr:
			count, _ := factory.CreateRepository().CountAlprs(*params)
			ctx.JSON(http.StatusOK, count)
			break
		default:
			ctx.JSON(http.StatusOK, 0)
		}
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

		switch p.AiType {
		case models.Od:
			err := factory.CreateRepository().DeleteOds(options)
			if err != nil {
				ctx.JSON(http.StatusOK, false)
				return
			}
			break
		case models.Fr:
			err := factory.CreateRepository().DeleteFrs(options)
			if err != nil {
				ctx.JSON(http.StatusOK, false)
				return
			}
			break
		case models.Alpr:
			err := factory.CreateRepository().DeleteAlprs(options)
			if err != nil {
				ctx.JSON(http.StatusOK, false)
				return
			}
			break
		}
		ctx.JSON(http.StatusOK, true)
	})
}
