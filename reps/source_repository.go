package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type SourceRepository struct {
	Connection *redis.Client
}

var redisKeySources = "sources:"

func (r *SourceRepository) Save(model *models.SourceModel) (*models.SourceModel, error) {
	conn := r.Connection
	if len(model.Id) == 0 {
		model.Id = NewId()
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

func (r *SourceRepository) GetById(id string) (*models.SourceModel, error) {
	conn := r.Connection
	key := redisKeySources + id
	var p models.SourceModel
	err := conn.HGetAll(context.Background(), key).Scan(&p)
	if err != nil {
		log.Println("Error getting source from redis: ", err)
		return nil, err
	}
	return &p, err
	//conn := r.Connection
	//key := redisKeySources + id
	//values, err := conn.HGetAll(context.Background(), key).Result()
	//if err != nil {
	//	log.Println("Error getting source from redis: ", err)
	//	return nil, err
	//}
	//var p models.SourceModel
	//err = p.Map(values)
	//return &p, err
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
