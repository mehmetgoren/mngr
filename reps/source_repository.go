package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type SourceRepository struct {
	Connection *redis.Client
}

var redisKeySources = "sources:"

func (r *SourceRepository) Save(model *models.SourceModel) (*models.SourceModel, error) {
	conn := r.Connection
	if len(model.Id) == 0 {
		model.Id = utils.NewId()
	}
	_, err := conn.HSet(context.Background(), redisKeySources+model.Id, Map(model)).Result()
	if err != nil {
		log.Println("Error while adding source: ", err)
		return nil, err
	}
	return model, nil
}

func (r *SourceRepository) GetAll() ([]*models.SourceModel, error) {
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeySources+"*").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			emptyList := make([]*models.SourceModel, 0)
			conn.Set(context.Background(), redisKeySources, emptyList, 0)
			return emptyList, nil
		} else {
			log.Println("Error getting sources from redis: ", err)
			return nil, err
		}
	}
	list := make([]*models.SourceModel, 0, 5)
	for _, key := range keys {
		var p models.SourceModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting source from redis: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *SourceRepository) Get(id string) (*models.SourceModel, error) {
	conn := r.Connection
	key := redisKeySources + id
	var p models.SourceModel
	err := conn.HGetAll(context.Background(), key).Scan(&p)
	if err != nil {
		log.Println("Error getting source from redis: ", err)
		return nil, err
	}
	return &p, err
}

func (r *SourceRepository) RemoveById(id string) error {
	conn := r.Connection
	_, err := conn.Del(context.Background(), redisKeySources+id).Result()
	if err != nil {
		log.Println("Error while deleting source: ", err)
		return err
	}
	return nil
}

func (r *SourceRepository) GetSourceStreamStatus(streamRepository *StreamRepository) ([]*models.SourceStatusModel, error) {
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeySources+"*").Result()
	if err != nil {
		log.Println("Error occurred while getting source status, err: ", err)
		emptyList := make([]*models.SourceStatusModel, 0)
		return emptyList, err
	}
	list := make([]*models.SourceStatusModel, 0, 5)
	for _, key := range keys {
		var source models.SourceModel
		err := conn.HGetAll(context.Background(), key).Scan(&source)
		if err != nil {
			log.Println("Error getting source from redis: ", err)
			return nil, err
		}
		sourceStatus := models.SourceStatusModel{SourceId: source.Id, Recording: false}
		stream, _ := streamRepository.Get(source.Id)
		sourceStatus.Streaming = stream != nil && len(stream.Id) > 0
		if sourceStatus.Streaming {
			sourceStatus.Recording = stream.RecordEnabled
		}

		list = append(list, &sourceStatus)
	}
	return list, nil
}
