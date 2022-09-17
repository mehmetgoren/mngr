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

func (r *RtspTemplateRepository) AddAll() (int64, error) {
	var total int64 = 0

	concord := &models.RtspTemplateModel{}
	concord.Id = "ConcordIpc"
	concord.Name = "ConcordIpc"
	concord.Description = "Generic Smart IP Camera"
	concord.Brand = "Concord"
	concord.DefaultUser = "admin"
	concord.DefaultPassword = "admin123456"
	concord.DefaultPort = "8554"
	concord.Route = "/profile0"
	concord.Templates = "rtsp://{user}:{password}@{ip}:{port}/{route}"
	temp, err := r.Connection.HSet(context.Background(), redisKeyRtspTemplate+concord.Id, Map(concord)).Result()
	total += temp

	tapo := &models.RtspTemplateModel{}
	tapo.Id = "Tapo"
	tapo.Name = "Tapo"
	tapo.Description = "TP-Lik Smart IP Camera Series"
	tapo.Brand = "TP-Link"
	tapo.DefaultUser = ""
	tapo.DefaultPassword = ""
	tapo.DefaultPort = "554"
	tapo.Route = "/stream1"
	tapo.Templates = "rtsp://{user}:{password}@{ip}:{port}/{route}"
	temp, err = r.Connection.HSet(context.Background(), redisKeyRtspTemplate+tapo.Id, Map(tapo)).Result()
	total += temp

	dahuaDvr := &models.RtspTemplateModel{}
	dahuaDvr.Id = "Dahua DVR"
	dahuaDvr.Name = "Dahua DVR"
	dahuaDvr.Description = "Dahua Analog Camera Series"
	dahuaDvr.Brand = "Dahua Ltd"
	dahuaDvr.DefaultUser = "admin"
	dahuaDvr.DefaultPassword = ""
	dahuaDvr.DefaultPort = "554"
	dahuaDvr.Route = "/cam/realmonitor?channel={camera_no}&subtype={subtype}"
	dahuaDvr.Templates = "rtsp://{user}:{password}@{ip}:{port}/{route}"
	temp, err = r.Connection.HSet(context.Background(), redisKeyRtspTemplate+dahuaDvr.Id, Map(dahuaDvr)).Result()
	total += temp

	ankerEufy2k := &models.RtspTemplateModel{}
	ankerEufy2k.Id = "Anker Eufy Security 2K"
	ankerEufy2k.Name = "Anker Eufy Security 2K"
	ankerEufy2k.Description = "Anker Smart IP Camera Series"
	ankerEufy2k.Brand = "Anker Ltd"
	ankerEufy2k.DefaultUser = ""
	ankerEufy2k.DefaultPassword = ""
	ankerEufy2k.DefaultPort = "554"
	ankerEufy2k.Route = "/live0"
	ankerEufy2k.Templates = "rtsp://{user}:{password}@{ip}:{port}/{route}"
	temp, err = r.Connection.HSet(context.Background(), redisKeyRtspTemplate+ankerEufy2k.Id, Map(ankerEufy2k)).Result()
	total += temp

	hikVision := &models.RtspTemplateModel{}
	hikVision.Id = "Hikvision IP Cameras"
	hikVision.Name = "Hikvision IP Cameras"
	hikVision.Description = "HIKVISION Leading IP Camera Series"
	hikVision.Brand = "HIKVISION Ltd"
	hikVision.DefaultUser = "admin"
	hikVision.DefaultPassword = ""
	hikVision.DefaultPort = "554"
	hikVision.Route = "/Streaming/Channels/101"
	hikVision.Templates = "rtsp://{user}:{password}@{ip}:{port}/{route}"
	temp, err = r.Connection.HSet(context.Background(), redisKeyRtspTemplate+hikVision.Id, Map(hikVision)).Result()
	total += temp

	return total, err
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
