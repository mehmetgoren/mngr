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
