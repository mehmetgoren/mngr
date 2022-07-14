package sqlt

import (
	"mngr/data"
	"mngr/utils"
	"strconv"
)

type SqliteRepository struct {
	Db *DbContext
}

func (s *SqliteRepository) GetOds(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.OdDto, error) {
	if ti == nil {
		return nil, nil
	}

	filter := &OdEntity{}
	filter.SourceId = sourceId
	filter.Year, _ = strconv.Atoi(ti.Year)
	filter.Month, _ = strconv.Atoi(ti.Month)
	filter.Day, _ = strconv.Atoi(ti.Day)
	filter.Hour, _ = strconv.Atoi(ti.Hour)
	entities := make([]*OdEntity, 0)
	q := s.Db.Ods.GetGormDb().Where(filter)
	if sort {
		q.Order("created_date desc")
	}
	result := q.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	if entities == nil || len(entities) == 0 {
		return nil, nil
	}

	mapper := &OdMapper{}
	ret := make([]*data.OdDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (s *SqliteRepository) GetFrs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.FrDto, error) {
	if ti == nil {
		return nil, nil
	}

	filter := &FrEntity{}
	filter.SourceId = sourceId
	filter.Year, _ = strconv.Atoi(ti.Year)
	filter.Month, _ = strconv.Atoi(ti.Month)
	filter.Day, _ = strconv.Atoi(ti.Day)
	filter.Hour, _ = strconv.Atoi(ti.Hour)
	entities := make([]*FrEntity, 0)
	q := s.Db.Frs.GetGormDb().Where(filter)
	if sort {
		q.Order("created_date desc")
	}
	result := q.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	if entities == nil || len(entities) == 0 {
		return nil, nil
	}

	mapper := &FrMapper{}
	ret := make([]*data.FrDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (s *SqliteRepository) GetAlprs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.AlprDto, error) {
	if ti == nil {
		return nil, nil
	}

	filter := &AlprEntity{}
	filter.SourceId = sourceId
	filter.Year, _ = strconv.Atoi(ti.Year)
	filter.Month, _ = strconv.Atoi(ti.Month)
	filter.Day, _ = strconv.Atoi(ti.Day)
	filter.Hour, _ = strconv.Atoi(ti.Hour)
	entities := make([]*AlprEntity, 0)
	q := s.Db.Alprs.GetGormDb().Where(filter)
	if sort {
		q.Order("created_date desc")
	}
	result := q.Find(&entities)
	if result.Error != nil {
		return nil, result.Error
	}
	if entities == nil || len(entities) == 0 {
		return nil, nil
	}

	mapper := &AlprMapper{}
	ret := make([]*data.AlprDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (s *SqliteRepository) RemoveOd(id string) error {
	//entity := &OdEntity{}
	//val, err := strconv.ParseUint(id, 10, 32)
	//if err != nil {
	//	return err
	//}
	//entity.ID = uint(val)
	//result := s.Db.Ods.GetGormDb().Unscoped().Model(entity).Delete(entity)
	result := s.Db.Ods.GetGormDb().Unscoped().Delete(&OdEntity{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
