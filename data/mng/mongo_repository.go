package mng

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mngr/data"
)

type MongoRepository struct {
	Db *DbContext
}

func createQuery(className string, params *data.GetParams) (bson.M, bson.D) {
	q := bson.M{}
	q["source_id"] = params.SourceId
	q["created_date"] = bson.M{"$gte": params.T1, "$lte": params.T2}
	if len(params.ClassName) > 0 {
		q[className] = primitive.Regex{Pattern: params.ClassName, Options: "i"} //bson.M{"$regex": "/.*" + params.ClassName + ".*/", "$options": "i"}
	}
	if params.NoPreparingVideoFile {
		q["video_file.name"] = bson.M{"$exists": true, "$ne": ""}
	}
	var sorts bson.D
	if params.Sort {
		sorts = bson.D{{"created_date", -1}}
	}
	return q, sorts
}

func (m *MongoRepository) GetOds(params *data.GetParams) ([]*data.OdDto, error) {
	if params == nil {
		return nil, nil
	}

	q, s := createQuery("detected_object.pred_cls_name", params)
	entities, err := m.Db.Ods.GetByQuery(q, s)
	if err != nil {
		return nil, err
	}
	if entities == nil && len(entities) == 0 {
		return nil, nil
	}

	mapper := &OdMapper{Config: m.Db.Config}
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

	q, s := createQuery("detected_face.pred_cls_name", params)
	entities, err := m.Db.Frs.GetByQuery(q, s)
	if err != nil {
		return nil, err
	}

	mapper := &FrMapper{Config: m.Db.Config}
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

	q, s := createQuery("detected_plate.plate", params)
	entities, err := m.Db.Alprs.GetByQuery(q, s)
	if err != nil {
		return nil, err
	}

	mapper := &AlprMapper{Config: m.Db.Config}
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
