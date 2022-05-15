package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type NetworkDiscoveryRepository struct {
	Connection *redis.Client
}

func getNetworkKey() string {
	return "onvif:network"
}

func (n *NetworkDiscoveryRepository) GetAll() (*models.NetworkDiscoveryModel, error) {
	conn := n.Connection
	key := getNetworkKey()
	js, err := conn.Get(context.Background(), key).Result()
	if err != nil {
		log.Println("Error getting network models from redis: ", err)
		return nil, err
	}
	if len(js) == 0 {
		return nil, nil
	}

	ret := &models.NetworkDiscoveryModel{}
	err = json.Unmarshal([]byte(js), &ret)
	if err != nil {
		log.Println("Error getting network models from redis: ", err)
		return nil, err
	}

	return ret, nil
}
