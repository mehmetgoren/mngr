package sqlt

import (
	"gorm.io/gorm"
	"log"
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
	if len(params.Module) > 0 {
		qarr = append(qarr, "module = ?")
		ps = append(ps, params.Module)
	}
	qarr = append(qarr, "created_date >= ? AND created_date < ?")
	ps = append(ps, params.T1, params.T2)

	if len(params.Label) > 0 {
		qarr = append(qarr, "label LIKE ?")
		ps = append(ps, "%"+params.Label+"%")
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

func (s *SqliteRepository) QueryAis(params data.QueryParams) ([]*data.AiDto, error) {
	ods := make([]*AiEntity, 0)
	db := s.Db.Ais.GetGormDb()
	db = createQuery(db, &params)
	result := db.Find(&ods)
	ret := make([]*data.AiDto, 0)
	if result.Error != nil {
		return ret, result.Error
	}

	mapper := &AiMapper{Config: s.Db.Config}
	for _, entity := range ods {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}
func (s *SqliteRepository) CountAis(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	db := s.Db.Ais.GetGormDb()
	db = createQuery(db, &params)
	return count[AiEntity](db)
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
				log.Println("an error occurred while deleting an ai image, err: " + err.Error())
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

func (s *SqliteRepository) DeleteAis(options *data.DeleteOptions) error {
	cast := func(t *AiEntity) *BaseEntity {
		return &t.BaseEntity
	}
	return deleteRec[AiEntity](s.Db.Ais.GetGormDb(), options, cast)
}

func (s *SqliteRepository) RemoveAi(id string) error {
	result := s.Db.Ais.GetGormDb().Unscoped().Delete(&AiEntity{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
