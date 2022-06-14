package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type RtspTemplateRepository struct {
	Connection *redis.Client
}

var redisKeyRtspTemplate = "rtsp_template:"

func (r *RtspTemplateRepository) GetAll() ([]*models.RtspTemplateModel, error) {
	ret := make([]*models.RtspTemplateModel, 0)
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeyRtspTemplate+"*").Result()
	if err != nil {
		return ret, err
	}
	for _, key := range keys {
		var p models.RtspTemplateModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting rtsp templates from redis: ", err)
			return nil, err
		}
		ret = append(ret, &p)
	}
	return ret, nil
}

type FailedStreamRepository struct {
	Connection *redis.Client
}

var redisKeyFailedStream = "failed_streams:"

func (f *FailedStreamRepository) GetAll() ([]*models.FailedStreamModel, error) {
	ret := make([]*models.FailedStreamModel, 0)
	conn := f.Connection
	keys, err := conn.Keys(context.Background(), redisKeyFailedStream+"*").Result()
	if err != nil {
		return ret, err
	}
	for _, key := range keys {
		var p models.FailedStreamModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting failed streams from redis: ", err)
			return nil, err
		}
		ret = append(ret, &p)
	}
	return ret, nil
}

type RecStuckRepository struct {
	Connection *redis.Client
}

var redisKeyRecStuck = "recstucks:"

func (r *RecStuckRepository) GetAll() ([]*models.RecStuckModel, error) {
	ret := make([]*models.RecStuckModel, 0)
	conn := r.Connection
	keys, err := conn.Keys(context.Background(), redisKeyRecStuck+"*").Result()
	if err != nil {
		return ret, err
	}
	for _, key := range keys {
		var p models.RecStuckModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting rec stuck from redis: ", err)
			return nil, err
		}
		ret = append(ret, &p)
	}
	return ret, nil
}
