package data

import (
	"mngr/models"
	"mngr/utils"
	"time"
)

type QueryParams struct {
	SourceId             string
	Module               string
	Label                string
	NoPreparingVideoFile bool
	T1                   time.Time
	T2                   time.Time
	Sort                 models.SortInfo
	Paging               models.PagingInfo
}

func GetParamsByHour(sourceId string, module string, dateStr string, sort models.SortInfo) *QueryParams {
	params := &QueryParams{}
	params.SourceId = sourceId
	params.Module = module
	params.Sort = sort
	params.SetupHourlyTimes(dateStr)
	return params
}

func (q *QueryParams) SetupHourlyTimes(dataStr string) {
	q.T1 = utils.StringToTime(dataStr)
	q.T2 = q.T1.Add(time.Hour)
}

func (q *QueryParams) SetupTimes(startDateStr string, endDateStr string) {
	q.T1 = utils.StringToTime(startDateStr)
	q.T2 = utils.StringToTime(endDateStr)
}

type DeleteOptions struct {
	Id          string
	DeleteImage bool
	DeleteVideo bool
}

type Repository interface {
	QueryAis(params QueryParams) ([]*AiDto, error)
	CountAis(params QueryParams) (int64, error)
	RemoveAi(id string) error
	DeleteAis(options *DeleteOptions) error
}
