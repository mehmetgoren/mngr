package data

import "mngr/utils"

type Repository interface {
	GetOds(sourceId string, ti *utils.TimeIndex, sort bool) ([]*OdDto, error)
	GetFrs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*FrDto, error)
	GetAlprs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*AlprDto, error)

	RemoveOd(id string) error
}
