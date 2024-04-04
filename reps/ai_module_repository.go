package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type AiModuleRepository struct {
	Connection *redis.Client
}

func getAiModuleKey(id string) string {
	return "ai_modules:" + id
}

func (r *AiModuleRepository) GetAiModules() ([]*models.AiModuleModel, error) {
	aiModules := make([]*models.AiModuleModel, 0)
	conn := r.Connection
	ctx := context.Background()
	keys, err := conn.Keys(ctx, getAiModuleKey("*")).Result()
	if err != nil {
		log.Println("Error getting ai modules from redis: ", err)
		return aiModules, err
	}

	for _, key := range keys {
		var aiModule models.AiModuleModel
		err := conn.HGetAll(ctx, key).Scan(&aiModule)
		if err != nil {
			log.Println("Error getting ai module from redis: ", err)
			return aiModules, err
		}
		aiModules = append(aiModules, &aiModule)
	}

	return aiModules, nil
}

func (r *AiModuleRepository) SaveAiModule(aiModule *models.AiModuleModel) (int, error) {
	conn := r.Connection
	if aiModule == nil {
		return 0, nil
	}
	_, err := conn.HSet(context.Background(), getAiModuleKey(aiModule.Name), Map(aiModule)).Result()
	if err != nil {
		log.Println("Error while adding a rtsp template: ", err.Error())
		return 0, err
	}
	return 1, nil
}

func (r *AiModuleRepository) RemoveAiModuleByName(name string) (int, error) {
	conn := r.Connection
	key := getAiModuleKey(name)
	err := conn.Del(context.Background(), key).Err()
	if err != nil {
		log.Println("Error while removing ai module: ", err)
		return 0, err
	}
	return 1, nil
}
