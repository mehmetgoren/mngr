package data

import (
	"mngr/utils"
	"time"
)

type GetParams struct {
	SourceId             string
	ClassName            string
	NoPreparingVideoFile bool
	T1                   time.Time
	T2                   time.Time
	Sort                 bool
}

func GetParamsByHour(sourceId string, dateStr string) *GetParams {
	params := &GetParams{}
	params.SourceId = sourceId
	params.Sort = true
	params.SetupTimes(dateStr)
	return params
}

func (g *GetParams) SetupTimes(dataStr string) {
	g.T1 = utils.StringToTime(dataStr)
	g.T2 = g.T1.Add(time.Hour)
}

type Repository interface {
	GetOds(params *GetParams) ([]*OdDto, error)
	GetFrs(params *GetParams) ([]*FrDto, error)
	GetAlprs(params *GetParams) ([]*AlprDto, error)

	RemoveOd(id string) error
}
