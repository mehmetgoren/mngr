package data

import (
	"mngr/models"
	"mngr/utils"
	"time"
)

type QueryParams struct {
	SourceId             string
	ClassName            string
	NoPreparingVideoFile bool
	T1                   time.Time
	T2                   time.Time
	Sort                 models.SortInfo
	Paging               models.PagingInfo
}

func GetParamsByHour(sourceId string, dateStr string, sort models.SortInfo) *QueryParams {
	params := &QueryParams{}
	params.SourceId = sourceId
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
	QueryOds(params QueryParams) ([]*OdDto, error)
	CountOds(params QueryParams) (int64, error)

	QueryFrs(params QueryParams) ([]*FrDto, error)
	CountFrs(params QueryParams) (int64, error)

	QueryAlprs(params QueryParams) ([]*AlprDto, error)
	CountAlprs(params QueryParams) (int64, error)

	RemoveOd(id string) error
	RemoveFr(id string) error
	RemoveAlpr(id string) error

	DeleteOds(options *DeleteOptions) error
	DeleteFrs(options *DeleteOptions) error
	DeleteAlprs(options *DeleteOptions) error
}
