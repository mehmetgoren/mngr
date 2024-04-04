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

func createQuery(label string, params *data.QueryParams) (bson.M, *PagingOptions, bson.D) {
	q := bson.M{}
	if len(params.SourceId) > 0 {
		q["source_id"] = params.SourceId
	}
	if len(params.Module) > 0 {
		q["module"] = params.Module
	}
	q["created_date"] = bson.M{"$gte": params.T1, "$lte": params.T2}
	if len(params.Label) > 0 {
		q[label] = primitive.Regex{Pattern: params.Label, Options: "i"}
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

func (m *MongoRepository) QueryAis(params data.QueryParams) ([]*data.AiDto, error) {
	if params.Sort.Enabled {
		if params.Sort.Field == "label" {
			params.Sort.Field = "detected_object.label"
		}
		if params.Sort.Field == "score" {
			params.Sort.Field = "detected_object.score"
		}
	}

	q, p, s := createQuery("detected_object.label", &params)
	entities, err := m.Db.Ais.GetByQuery(q, p, s)
	if err != nil {
		return nil, err
	}
	if entities == nil && len(entities) == 0 {
		return nil, nil
	}

	mapper := &AiMapper{Config: m.Db.Config}
	ret := make([]*data.AiDto, 0)
	for _, entity := range entities {
		ret = append(ret, mapper.Map(entity))
	}

	return ret, nil
}
func (m *MongoRepository) CountAis(params data.QueryParams) (int64, error) {
	params.Sort.Enabled = false
	params.Paging.Enabled = false
	f, _, _ := createQuery("detected_object.label", &params)
	return count(m.Db.Ais.coll, f)
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
			log.Println("an error occurred while deleting an ai image, err: " + err.Error())
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

func (m *MongoRepository) DeleteAis(options *data.DeleteOptions) error {
	getGroupId := func(t *AiEntity) string {
		return t.GroupId
	}
	getImageFileName := func(t *AiEntity) string {
		return t.ImageFileName
	}
	getVideoFileName := func(t *AiEntity) string {
		return t.VideoFile.Name
	}
	return deleteRec[AiEntity](m.Db.Ais.GetCollection(), options, getGroupId, getImageFileName, getVideoFileName)
}

func (m *MongoRepository) RemoveAi(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = m.Db.Ais.DeleteOneById(objectId)

	return err
}
