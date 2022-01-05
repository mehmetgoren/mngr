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

func (r *SourceRepository) GetAllSources() ([]*models.Source, error) {
	conn := r.Connection
	//jsonList, err := conn.Get(context.Background(), redisKeySources).Result()
	keys, err := conn.Keys(context.Background(), redisKeySources+"*").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			emptyList := make([]*models.Source, 0)
			conn.Set(context.Background(), redisKeySources, emptyList, 0)
			return emptyList, nil
		} else {
			log.Println("Error getting sources from redis: ", err)
			return nil, err
		}
	}
	var list []*models.Source
	for _, key := range keys {
		dic, _ := conn.HGetAll(context.Background(), key).Result()
		var p models.Source
		p.TypeName = dic["type_name"]
		p.Name = dic["name"]
		p.Brand = dic["brand"]
		p.RtspAddress = dic["rtsp_address"]
		p.Id = dic["id"]
		list = append(list, &p)
	}
	return list, nil
}
