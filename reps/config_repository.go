package reps

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
)

type ConfigRepository struct {
	Connection *redis.Client
}

func (r *ConfigRepository) GetConfig() (*models.Config, error) {
	var config = &models.Config{}
	conn := r.Connection
	key := "config"
	data, err := conn.Get(context.Background(), key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return config, nil
		} else {
			log.Println("Error getting sources from redis: ", err)
			return nil, err
		}
	}

	err = json.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *ConfigRepository) SaveConfig(config *models.Config) error {
	conn := r.Connection
	key := "config"
	data, err := json.Marshal(config)
	if err != nil {
		log.Println("Error serializing json: ", err)
		return err
	}

	err = conn.Set(context.Background(), key, data, 0).Err()
	if err != nil {
		log.Println("Error saving ml config to redis: ", err)
		return err
	}

	return nil
}

func (r *ConfigRepository) RestoreConfig() (*models.Config, error) {
	//todo: replace this json with the new one
	objJson := `{"device":{"device_name":"dev","device_arch":0},"general":{"dir_paths":[]},"db":{"type":1,"connection_string":"mongodb://localhost:27017"},"ffmpeg":{"use_double_quotes_for_path":false,"max_operation_retry_count":10000000,"ms_init_interval":3,"watch_dog_interval":23,"watch_dog_failed_wait_interval":3,"start_task_wait_for_interval":1,"record_concat_limit":1,"record_video_file_indexer_interval":60,"ms_port_start":7000,"ms_port_end":8000},"ai":{"video_clip_duration":10},"sense_ai":{"image":2,"host":"127.0.0.1","port":32168},"jobs":{"mac_ip_matching_enabled":false,"mac_ip_matching_interval":120,"black_screen_monitor_enabled":false,"black_screen_monitor_interval":600},"archive":{"limit_percent":95,"action_type":0,"move_location":""},"snapshot":{"process_count":4,"overlay":true},"hub":{"enabled":false,"address":"http://localhost:5268","token":"","web_app_address":"http://localhost:8080","max_retry":100}}`

	conn := r.Connection
	key := "config"
	err := conn.Set(context.Background(), key, objJson, 0).Err()
	if err != nil {
		log.Println("Error saving default ml config to redis: ", err)
		return nil, err
	}

	var config = &models.Config{}
	err = json.Unmarshal([]byte(objJson), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
