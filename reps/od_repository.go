package reps

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type OdRepository struct {
	Connection *redis.Client
}

var redisKeyOds = "ods:"

func (o *OdRepository) Get(id string) (*models.OdModel, error) {
	conn := o.Connection
	key := redisKeyOds + id
	var p models.OdModel
	err := conn.HGetAll(context.Background(), key).Scan(&p)
	if err != nil {
		log.Println("Error getting object detection model from redis: ", err)
		return nil, err
	}
	return &p, err
}

func (o *OdRepository) Save(model *models.OdModel) (*models.OdModel, error) {
	conn := o.Connection
	if len(model.Id) == 0 {
		errMsg := "insert operation is not supported on mngr"
		log.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	_, err := conn.HSet(context.Background(), redisKeyOds+model.Id, Map(model)).Result()
	if err != nil {
		log.Println("Error while adding object detection model: ", err)
		return nil, err
	}
	return model, nil
}

func (o *OdRepository) RemoveById(id string) error {
	conn := o.Connection
	_, err := conn.Del(context.Background(), redisKeyOds+id).Result()
	if err != nil {
		log.Println("Error while deleting object detection model: ", err)
		return err
	}
	return nil
}
