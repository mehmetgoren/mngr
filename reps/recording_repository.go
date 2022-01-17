package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"mngr/models"
)

type RecordingRepository struct {
	Connection *redis.Client
}

var redisKeyRecording = "recording:"

func (r *RecordingRepository) Get(id string) (*models.Recording, error) {
	key := redisKeyRecording + id
	_, err := r.Connection.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	recording := &models.Recording{}
	err = r.Connection.HGetAll(context.Background(), redisKeyRecording+id).Scan(recording)
	if err != nil {
		return nil, err
	}
	return recording, nil
}
