package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"strconv"
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

type VariousInfosRepository struct {
	Connection *redis.Client
}

func (v *VariousInfosRepository) Get() (*models.VariousInfos, error) {
	conn := v.Connection
	ctx := context.Background()
	result, err := conn.HGet(ctx, "rtmpports", "ports_count").Result()
	if err != nil {
		return nil, err
	}
	ret := &models.VariousInfos{}
	ret.RtmpPortCounter, _ = strconv.Atoi(result)

	//result2, err := conn.LRange(ctx, "zombies:docker", 0, -1).Result()
	result2, err := conn.SMembers(ctx, "zombies:docker").Result()
	if err != nil {
		return nil, err
	}
	if result2 != nil {
		ret.RtmpContainerZombies = make([]string, 0)
		for _, v := range result2 {
			ret.RtmpContainerZombies = append(ret.RtmpContainerZombies, v)
		}
	}

	result3, err := conn.SMembers(ctx, "zombies:ffmpeg").Result()
	if err != nil {
		return nil, err
	}
	if result3 != nil {
		ret.FFmpegProcessZombies = make([]int, 0)
		for _, v := range result3 {
			iv, _ := strconv.Atoi(v)
			ret.FFmpegProcessZombies = append(ret.FFmpegProcessZombies, iv)
		}
	}

	return ret, nil
}
