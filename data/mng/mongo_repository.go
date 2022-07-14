package mng

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mngr/data"
	"mngr/utils"
)

type MongoRepository struct {
	Db *DbContext
}

func (m *MongoRepository) GetOds(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.OdDto, error) {
	if ti == nil {
		return nil, nil
	}

	var sorts bson.D
	if sort {
		sorts = bson.D{{"created_date", -1}}
	}
	year, month, day, hour := ti.GetValuesAsInt()
	entities, err := m.Db.Ods.GetByQuery(bson.M{"source_id": sourceId, "year": year, "month": month, "day": day, "hour": hour}, sorts)
	if err != nil {
		return nil, err
	}
	if entities == nil && len(entities) == 0 {
		return nil, nil
	}

	mapper := &OdMapper{}
	ret := make([]*data.OdDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (m *MongoRepository) GetFrs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.FrDto, error) {
	if ti == nil {
		return nil, nil
	}

	var sorts bson.D
	if sort {
		sorts = bson.D{{"created_date", -1}}
	}
	year, month, day, hour := ti.GetValuesAsInt()
	entities, err := m.Db.Frs.GetByQuery(bson.M{"source_id": sourceId, "year": year, "month": month, "day": day, "hour": hour}, sorts)
	if err != nil {
		return nil, err
	}

	mapper := &FrMapper{}
	ret := make([]*data.FrDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (m *MongoRepository) GetAlprs(sourceId string, ti *utils.TimeIndex, sort bool) ([]*data.AlprDto, error) {
	if ti == nil {
		return nil, nil
	}

	var sorts bson.D
	if sort {
		sorts = bson.D{{"created_date", -1}}
	}
	year, month, day, hour := ti.GetValuesAsInt()
	entities, err := m.Db.Alprs.GetByQuery(bson.M{"source_id": sourceId, "year": year, "month": month, "day": day, "hour": hour}, sorts)
	if err != nil {
		return nil, err
	}

	mapper := &AlprMapper{}
	ret := make([]*data.AlprDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}

func (m *MongoRepository) RemoveOd(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Ods.DeleteOneById(objectId)

	return err
}
