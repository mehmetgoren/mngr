package sqlt

import (
	"mngr/data"
)

type SqliteRepository struct {
	Db *DbContext
}

var qs = "source_id = ? AND created_date >= ? AND created_date < ?"

func (s *SqliteRepository) GetOds(params *data.GetParams) ([]*data.OdDto, error) {
	if params == nil {
		return nil, nil
	}

	ods := make([]*OdEntity, 0)
	db := s.Db.Ods.GetGormDb()
	q := db.Where(qs, params.SourceId, params.T1, params.T2)
	if params.Sort {
		q.Order("created_date desc")
	}
	result := q.Find(&ods)
	ret := make([]*data.OdDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &OdMapper{}
	for _, entity := range ods {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (s *SqliteRepository) GetFrs(params *data.GetParams) ([]*data.FrDto, error) {
	if params == nil {
		return nil, nil
	}

	frs := make([]*FrEntity, 0)
	db := s.Db.Frs.GetGormDb()
	q := db.Where(qs, params.SourceId, params.T1, params.T2)
	if params.Sort {
		q.Order("created_date desc")
	}
	result := q.Find(&frs)
	ret := make([]*data.FrDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &FrMapper{}
	for _, fr := range frs {
		ret = append(ret, mapper.Map(fr))
	}

	return ret, nil
}

func (s *SqliteRepository) GetAlprs(params *data.GetParams) ([]*data.AlprDto, error) {
	if params == nil {
		return nil, nil
	}

	alprs := make([]*AlprEntity, 0)
	db := s.Db.Alprs.GetGormDb()
	q := db.Where(qs, params.SourceId, params.T1, params.T2)
	if params.Sort {
		q.Order("created_date desc")
	}
	result := q.Find(&alprs)
	ret := make([]*data.AlprDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &AlprMapper{}
	for _, alpr := range alprs {
		ret = append(ret, mapper.Map(alpr))
	}

	return ret, nil
}

func (s *SqliteRepository) RemoveOd(id string) error {
	result := s.Db.Ods.GetGormDb().Unscoped().Delete(&OdEntity{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
