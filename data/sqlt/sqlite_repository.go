package sqlt

import (
	"gorm.io/gorm"
	"mngr/data"
	"strings"
)

type SqliteRepository struct {
	Db *DbContext
}

func createQuery(db *gorm.DB, params *data.GetParams) *gorm.DB {
	var qs strings.Builder
	ps := make([]interface{}, 0)
	qs.WriteString("source_id = ? AND created_date >= ? AND created_date < ?")
	ps = append(ps, params.SourceId, params.T1, params.T2)
	var q *gorm.DB
	if len(params.ClassName) > 0 {
		qs.WriteString(" AND pred_cls_name like ?")
		ps = append(ps, "%"+params.ClassName+"%")

	}
	if params.NoPreparingVideoFile {
		qs.WriteString(" AND video_file_name is not null and LENGTH(video_file_name) > 0")
	}
	q = db.Where(qs.String(), ps...)
	if params.Sort {
		q.Order("created_date desc")
	}
	return q
}

func (s *SqliteRepository) GetOds(params *data.GetParams) ([]*data.OdDto, error) {
	if params == nil {
		return nil, nil
	}

	ods := make([]*OdEntity, 0)
	db := s.Db.Ods.GetGormDb()
	db = createQuery(db, params)
	result := db.Find(&ods)
	ret := make([]*data.OdDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &OdMapper{Config: s.Db.Config}
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
	db = createQuery(db, params)
	result := db.Find(&frs)
	ret := make([]*data.FrDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &FrMapper{Config: s.Db.Config}
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
	db = createQuery(db, params)
	result := db.Find(&alprs)
	ret := make([]*data.AlprDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &AlprMapper{Config: s.Db.Config}
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
