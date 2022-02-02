package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type StreamingRepository struct {
	Connection *redis.Client
}

var redisKeyStreaming = "streaming:"

func (r *StreamingRepository) Get(id string) (*models.StreamingModel, error) {
	conn := r.Connection
	key := redisKeyStreaming + id
	var model models.StreamingModel
	err := conn.HGetAll(context.Background(), key).Scan(&model)
	if err != nil {
		log.Println("Error getting streaming from redis: ", err)
		return nil, err
	}

	return &model, nil
}

func (r *StreamingRepository) GetAll() ([]*models.StreamingModel, error) {
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeyStreaming+"*").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			emptyList := make([]*models.StreamingModel, 0)
			conn.Set(context.Background(), redisKeyStreaming, emptyList, 0)
			return emptyList, nil
		} else {
			log.Println("Error getting all streaming from redis: ", err)
			return nil, err
		}
	}
	list := make([]*models.StreamingModel, 0, 5)
	for _, key := range keys {
		var p models.StreamingModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting streaming from redis: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}
