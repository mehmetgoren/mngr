package sqlt

import (
	"gorm.io/gorm"
	"mngr/data"
	"mngr/models"
	"os"
	"strconv"
	"strings"
)

type SqliteRepository struct {
	Db *DbContext
}

func count[T any](db *gorm.DB) (int64, error) {
	t := new(T)
	var count int64
	result := db.Model(t).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func createQuery(db *gorm.DB, params *data.QueryParams) *gorm.DB {
	var qarr = make([]string, 0)
	ps := make([]interface{}, 0)
	if len(params.SourceId) > 0 {
		qarr = append(qarr, "source_id = ?")
		ps = append(ps, params.SourceId)
	}
	qarr = append(qarr, "created_date >= ? AND created_date < ?")
	ps = append(ps, params.T1, params.T2)

	if len(params.ClassName) > 0 {
		qarr = append(qarr, "pred_cls_name LIKE ?")
		ps = append(ps, "%"+params.ClassName+"%")
	}
	if params.NoPreparingVideoFile {
		qarr = append(qarr, "video_file_name IS NOT NULL AND LENGTH(video_file_name) > 0")
	}
	qs := strings.Join(qarr, " AND ")

	var q = db
	if params.Paging.Enabled {
		offset := (params.Paging.Page - 1) * params.Paging.Take
		q = db.Offset(offset).Limit(params.Paging.Take)
	}
	if params.Sort.Enabled {
		var suffix string
		switch params.Sort.Sort {
		case models.Ascending:
			suffix = " ASC"
			break
		case models.Descending:
			suffix = " DESC"
			break
		}
		q = q.Order(params.Sort.Field + suffix)
	}

	q = q.Where(qs, ps...)
	return q
}

func (s *SqliteRepository) QueryOds(params data.QueryParams) ([]*data.OdDto, error) {
	ods := make([]*OdEntity, 0)
	db := s.Db.Ods.GetGormDb()
	db = createQuery(db, &params)
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
func (s *SqliteRepository) CountOds(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	db := s.Db.Ods.GetGormDb()
	db = createQuery(db, &params)
	return count[OdEntity](db)
}

func (s *SqliteRepository) QueryFrs(params data.QueryParams) ([]*data.FrDto, error) {
	frs := make([]*FrEntity, 0)
	db := s.Db.Frs.GetGormDb()
	db = createQuery(db, &params)
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
func (s *SqliteRepository) CountFrs(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	db := s.Db.Frs.GetGormDb()
	db = createQuery(db, &params)
	return count[FrEntity](db)
}

func (s *SqliteRepository) QueryAlprs(params data.QueryParams) ([]*data.AlprDto, error) {
	alprs := make([]*AlprEntity, 0)
	db := s.Db.Alprs.GetGormDb()
	db = createQuery(db, &params)
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
func (s *SqliteRepository) CountAlprs(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	db := s.Db.Alprs.GetGormDb()
	db = createQuery(db, &params)
	return count[AlprEntity](db)
}

func (s *SqliteRepository) RemoveOd(id string) error {
	result := s.Db.Ods.GetGormDb().Unscoped().Delete(&OdEntity{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func deleteRec[T any](db *gorm.DB, options *data.DeleteOptions, cast func(t *T) *BaseEntity) error {
	if options.Id == "" {
		return nil
	}

	id_, err := strconv.ParseUint(options.Id, 10, 32)
	if err != nil {
		return err
	}
	id := uint(id_)
	entity := new(T)
	result := db.First(entity, id)
	if result.Error != nil {
		return result.Error
	}

	result = db.Delete(entity, id)
	if result.Error != nil {
		return result.Error
	}

	if options.DeleteImage {
		be := cast(entity)
		if be.GroupId != "" {
			result = db.Where("group_id = ?", be.GroupId).Delete(entity)
			if result.Error != nil {
				return result.Error
			}
			err = os.Remove(be.ImageFileName)
			if err != nil {
				return err
			}
		}
	}

	if options.DeleteVideo {
		be := cast(entity)
		if be.VideoFileName != "" {
			result = db.Where("video_file_name = ?", be.VideoFileName).Delete(entity)
			if result.Error != nil {
				return result.Error
			}
			err = os.Remove(be.VideoFileName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *SqliteRepository) DeleteOds(options *data.DeleteOptions) error {
	cast := func(t *OdEntity) *BaseEntity {
		return &t.BaseEntity
	}
	return deleteRec[OdEntity](s.Db.Ods.GetGormDb(), options, cast)
}

func (s *SqliteRepository) DeleteFrs(options *data.DeleteOptions) error {
	cast := func(t *FrEntity) *BaseEntity {
		return &t.BaseEntity
	}
	return deleteRec[FrEntity](s.Db.Frs.GetGormDb(), options, cast)
}

func (s *SqliteRepository) DeleteAlprs(options *data.DeleteOptions) error {
	cast := func(t *AlprEntity) *BaseEntity {
		return &t.BaseEntity
	}
	return deleteRec[AlprEntity](s.Db.Alprs.GetGormDb(), options, cast)
}
