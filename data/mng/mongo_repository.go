package mng

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mngr/data"
)

type MongoRepository struct {
	Db *DbContext
}

func createQuery(params *data.GetParams) bson.M {
	return bson.M{"source_id": params.SourceId, "created_date": bson.M{"$gte": params.T1, "$lt": params.T2}}
}

func (m *MongoRepository) GetOds(params *data.GetParams) ([]*data.OdDto, error) {
	if params == nil {
		return nil, nil
	}

	var sorts bson.D
	if params.Sort {
		sorts = bson.D{{"created_date", -1}}
	}
	entities, err := m.Db.Ods.GetByQuery(createQuery(params), sorts)
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

func (m *MongoRepository) GetFrs(params *data.GetParams) ([]*data.FrDto, error) {
	if params == nil {
		return nil, nil
	}

	var sorts bson.D
	if params.Sort {
		sorts = bson.D{{"created_date", -1}}
	}
	entities, err := m.Db.Frs.GetByQuery(createQuery(params), sorts)
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

func (m *MongoRepository) GetAlprs(params *data.GetParams) ([]*data.AlprDto, error) {
	if params == nil {
		return nil, nil
	}

	var sorts bson.D
	if params.Sort {
		sorts = bson.D{{"created_date", -1}}
	}
	entities, err := m.Db.Alprs.GetByQuery(createQuery(params), sorts)
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
