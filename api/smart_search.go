package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"net/http"
	"strings"
)

func RegisterSmartSearchEndpoints(router *gin.Engine, factory *cmn.Factory) {
	router.POST("smartsearch", func(ctx *gin.Context) {
		var p models.SmartSearchParams
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusOK, emptyList)
			return
		}

		mapFn := func(par *models.SmartSearchParams) *data.QueryParams {
			ret := &data.QueryParams{}
			ret.SourceId = p.SourceId
			ret.SetupTimes(p.StartDateTimeStr, p.EndDateTimeStr)
			ret.ClassName = p.PredClassName
			ret.Sort.Enabled = false
			ret.Paging.Enabled = false
			return ret
		}

		ret := make([]*data.AiDataDto, 0)
		params := mapFn(&p)
		ods, _ := factory.CreateRepository().QueryOds(*params)
		if ods != nil && len(ods) > 0 {
			if strings.HasPrefix(p.Color.HexColor, "#") {
				p.Color.HexColor = p.Color.HexColor[1:]
			}
			hex := models.Hex(p.Color.HexColor)
			rgb, _ := models.Hex2RGB(hex)

			c := colorful.Color{R: float64(rgb.Red), G: float64(rgb.Green), B: float64(rgb.Blue)}

			for _, od := range ods {
				metadata := od.DetectedObject.Metadata
				if metadata == nil {
					continue
				}
				colors := metadata.Colors
				if colors == nil || len(colors) == 0 {
					continue
				}
				thresh := p.Color.Threshold
				for _, color := range colors {
					c2 := colorful.Color{R: float64(color.R), G: float64(color.G), B: float64(color.B)}
					distance := getDistanceMethod(&p.Color, &c)(c2)
					if distance <= thresh {
						ai := &data.AiDataDto{}
						ai.MapFromOd(od)
						ret = append(ret, ai)
						break
					}
				}
			}
		}

		ctx.JSON(http.StatusOK, ret)
	})
}

func getDistanceMethod(p *models.SmartSearchColor, c *colorful.Color) func(c2 colorful.Color) float64 {
	switch p.DifferenceMethod {
	case "CIE76":
		return c.DistanceCIE76
	case "CIE94":
		return c.DistanceCIE94
	case "Rgb":
		return c.DistanceRgb
	case "Lab":
		return c.DistanceLab
	case "Luv":
		return c.DistanceLuv
	case "LinearRGB":
		return c.DistanceLinearRGB
	case "HPLuv":
		return c.DistanceHPLuv
	case "HSLuv":
		return c.DistanceHSLuv
	default: // CIEDE2000
		return c.DistanceCIEDE2000
	}
}
