package api

import (
	"github.com/gin-gonic/gin"
	"mngr/data"
	"mngr/data/cmn"
	"net/http"
)

const (
	Od   = 0
	Fr   = 1
	Alpr = 2
)

type QueryAiDataParams struct {
	AiType               int    `json:"ai_type"`
	SourceId             string `json:"source_id"`
	DateTimeStr          string `json:"date_time_str"`
	PredClassName        string `json:"pred_class_name"`
	NoPreparingVideoFile bool   `json:"no_preparing_video_file"`
}

func RegisterAiDataEndpoints(router *gin.Engine, factory *cmn.Factory) {
	router.POST("queryaidata", func(ctx *gin.Context) {
		ret := make([]*data.AiDataDto, 0)
		var p QueryAiDataParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, ret)
			return
		}

		params := &data.GetParams{}
		params.SourceId = p.SourceId
		params.ClassName = p.PredClassName
		params.Sort = true
		params.NoPreparingVideoFile = p.NoPreparingVideoFile
		params.SetupTimes(p.DateTimeStr)

		switch p.AiType {
		case Od:
			ods, _ := factory.CreateRepository().GetOds(params)
			for _, od := range ods {
				ai := &data.AiDataDto{}
				ai.MapFromOd(od)
				ret = append(ret, ai)
			}
			break
		case Fr:
			frs, _ := factory.CreateRepository().GetFrs(params)
			for _, fr := range frs {
				ai := &data.AiDataDto{}
				ai.MapFromFr(fr)
				ret = append(ret, ai)
			}
			break
		case Alpr:
			alprs, _ := factory.CreateRepository().GetAlprs(params)
			for _, alpr := range alprs {
				ai := &data.AiDataDto{}
				ai.MapFromAlpr(alpr)
				ret = append(ret, ai)
			}
			break
		}
		ctx.JSON(http.StatusOK, ret)
	})
}
