package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type OnvifRepository struct {
	Connection *redis.Client
}

func getTargetKey(address string) string {
	ip := utils.ParseIp(address)
	return "onvif:" + ip
}

func (o *OnvifRepository) Get(address string) (*models.OnvifModel, error) {
	conn := o.Connection
	key := getTargetKey(address)
	js, err := conn.Get(context.Background(), key).Result()
	if len(js) == 0 {
		return nil, nil
	}
	model := &models.OnvifModel{}
	err = json.Unmarshal([]byte(js), model)
	if err != nil {
		log.Println("Error getting target model from redis: ", err)
		return nil, err
	}

	return model, nil
}
