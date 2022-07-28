package mng

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mngr/data"
)

type MongoRepository struct {
	Db *DbContext
}

func count(coll *mongo.Collection, filter bson.M) (int64, error) {
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func createQuery(className string, params *data.QueryParams) (bson.M, *PagingOptions, bson.D) {
	q := bson.M{}
	if len(params.SourceId) > 0 {
		q["source_id"] = params.SourceId
	}
	q["created_date"] = bson.M{"$gte": params.T1, "$lte": params.T2}
	if len(params.ClassName) > 0 {
		q[className] = primitive.Regex{Pattern: params.ClassName, Options: "i"}
	}
	if params.NoPreparingVideoFile {
		q["video_file.name"] = bson.M{"$exists": true, "$ne": ""}
	}

	var po *PagingOptions = nil
	if params.Paging.Enabled {
		po = &PagingOptions{
			Skip:  params.Paging.Page*params.Paging.Take - params.Paging.Take,
			Limit: params.Paging.Take,
		}
	}

	var sorts bson.D
	if params.Sort.Enabled {
		sorts = bson.D{{params.Sort.Field, params.Sort.Sort}}
	}
	return q, po, sorts
}

func (m *MongoRepository) QueryOds(params data.QueryParams) ([]*data.OdDto, error) {
	if params.Sort.Enabled {
		if params.Sort.Field == "pred_cls_name" {
			params.Sort.Field = "detected_object.pred_cls_name"
		}
		if params.Sort.Field == "pred_score" {
			params.Sort.Field = "detected_object.pred_score"
		}
	}

	q, p, s := createQuery("detected_object.pred_cls_name", &params)
	entities, err := m.Db.Ods.GetByQuery(q, p, s)
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
func (m *MongoRepository) CountOds(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	f, _, _ := createQuery("detected_object.pred_cls_name", &params)
	return count(m.Db.Ods.coll, f)
}

func (m *MongoRepository) QueryFrs(params data.QueryParams) ([]*data.FrDto, error) {
	if params.Sort.Enabled {
		if params.Sort.Field == "pred_cls_name" {
			params.Sort.Field = "detected_face.pred_cls_name"
		}
		if params.Sort.Field == "pred_score" {
			params.Sort.Field = "detected_face.pred_score"
		}
	}

	q, p, s := createQuery("detected_face.pred_cls_name", &params)
	entities, err := m.Db.Frs.GetByQuery(q, p, s)
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
func (m *MongoRepository) CountFrs(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	q, _, _ := createQuery("detected_face.pred_cls_name", &params)
	return count(m.Db.Frs.coll, q)
}

func (m *MongoRepository) QueryAlprs(params data.QueryParams) ([]*data.AlprDto, error) {
	if params.Sort.Enabled {
		if params.Sort.Field == "pred_cls_name" {
			params.Sort.Field = "detected_plate.plate"
		}
		if params.Sort.Field == "pred_score" {
			params.Sort.Field = "detected_plate.confidence"
		}
	}
	q, p, s := createQuery("detected_plate.plate", &params)
	entities, err := m.Db.Alprs.GetByQuery(q, p, s)
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
func (m *MongoRepository) CountAlprs(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	q, _, _ := createQuery("detected_plate.plate", &params)
	return count(m.Db.Alprs.coll, q)
}

func (m *MongoRepository) RemoveOd(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Ods.DeleteOneById(objectId)

	return err
}
