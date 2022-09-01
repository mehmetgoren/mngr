package mng

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mngr/data"
	"os"
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

func deleteRec[T any](coll *mongo.Collection, options *data.DeleteOptions,
	getGroupId func(t *T) string, getImageFileName func(t *T) string, getVideoFileName func(t *T) string) error {
	objectId, err := primitive.ObjectIDFromHex(options.Id)
	if err != nil {
		log.Println("Invalid id")
	}

	ctx := context.TODO()
	result := coll.FindOne(ctx, bson.M{"_id": objectId})
	entity := new(T)
	err = result.Decode(entity)
	if err != nil {
		return err
	}
	_, err = coll.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	if options.DeleteImage {
		_, err = coll.DeleteMany(ctx, bson.M{"group_id": getGroupId(entity)})
		if err != nil {
			return err
		}
		imageFileName := getImageFileName(entity)
		err = os.Remove(imageFileName)
		if err != nil {
			return err
		}
	}

	if options.DeleteVideo {
		videoFileName := getVideoFileName(entity)
		_, err = coll.DeleteMany(ctx, bson.M{"video_file.name": videoFileName})
		if err != nil {
			return err
		}
		err = os.Remove(videoFileName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoRepository) DeleteOds(options *data.DeleteOptions) error {
	getGroupId := func(t *OdEntity) string {
		return t.GroupId
	}
	getImageFileName := func(t *OdEntity) string {
		return t.ImageFileName
	}
	getVideoFileName := func(t *OdEntity) string {
		return t.VideoFile.Name
	}
	return deleteRec[OdEntity](m.Db.Ods.GetCollection(), options, getGroupId, getImageFileName, getVideoFileName)
}

func (m *MongoRepository) DeleteFrs(options *data.DeleteOptions) error {
	getGroupId := func(t *FrEntity) string {
		return t.GroupId
	}
	getImageFileName := func(t *FrEntity) string {
		return t.ImageFileName
	}
	getVideoFileName := func(t *FrEntity) string {
		return t.VideoFile.Name
	}
	return deleteRec[FrEntity](m.Db.Frs.GetCollection(), options, getGroupId, getImageFileName, getVideoFileName)
}

func (m *MongoRepository) DeleteAlprs(options *data.DeleteOptions) error {
	getGroupId := func(t *AlprEntity) string {
		return t.GroupId
	}
	getImageFileName := func(t *AlprEntity) string {
		return t.ImageFileName
	}
	getVideoFileName := func(t *AlprEntity) string {
		return t.VideoFile.Name
	}
	return deleteRec[AlprEntity](m.Db.Alprs.GetCollection(), options, getGroupId, getImageFileName, getVideoFileName)
}

func (m *MongoRepository) RemoveOd(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Ods.DeleteOneById(objectId)

	return err
}

func (m *MongoRepository) RemoveFr(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Frs.DeleteOneById(objectId)

	return err
}

func (m *MongoRepository) RemoveAlpr(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Alprs.DeleteOneById(objectId)

	return err
}
