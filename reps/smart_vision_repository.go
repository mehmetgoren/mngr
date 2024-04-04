package reps

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type SmartVisionRepository struct {
	Connection *redis.Client
}

var redisKeySmartVisions = "smart_visions:"

func (o *SmartVisionRepository) Get(id string) (*models.SmartVisionModel, error) {
	conn := o.Connection
	key := redisKeySmartVisions + id
	var p models.SmartVisionModel
	err := conn.HGetAll(context.Background(), key).Scan(&p)
	if err != nil {
		log.Println("Error getting object detection model from redis: ", err)
		return nil, err
	}
	return &p, err
}

func (o *SmartVisionRepository) GetAll() ([]*models.SmartVisionModel, error) {
	ret := make([]*models.SmartVisionModel, 0)
	conn := o.Connection
	allKey := redisKeySmartVisions + "*"
	keys, err := conn.Keys(context.Background(), allKey).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return ret, nil
		} else {
			log.Println("Error getting all stream from redis: ", err)
			return nil, err
		}
	}
	for _, key := range keys {
		var p models.SmartVisionModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting object detection model from redis: ", err)
			return nil, err
		}
		ret = append(ret, &p)
	}

	return ret, err
}

func (o *SmartVisionRepository) Save(model *models.SmartVisionModel) (*models.SmartVisionModel, error) {
	conn := o.Connection
	if len(model.Id) == 0 {
		errMsg := "insert operation is not supported on mngr"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	_, err := conn.HSet(context.Background(), redisKeySmartVisions+model.Id, Map(model)).Result()
	if err != nil {
		log.Println("Error while adding object detection model: ", err)
		return nil, err
	}
	return model, nil
}

func (o *SmartVisionRepository) RemoveById(id string) error {
	conn := o.Connection
	_, err := conn.Del(context.Background(), redisKeySmartVisions+id).Result()
	if err != nil {
		log.Println("Error while deleting object detection model: ", err)
		return err
	}
	return nil
}
