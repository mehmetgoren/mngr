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
	objJson := `{"device":{"device_name":"dev","device_type":0},"jetson":{"model_name":"ssd-mobilenet-v2"},"torch":{"model_name":"ultralytics/yolov5","model_name_specific":"yolov5x6"},"tensorflow":{"model_name":"efficientdet/lite4/detection","cache_folder":"/mnt/sdc1/test_projects/tf_cache"},"source_reader":{"resize_img":false,"buffer_size":2,"max_retry":150,"max_retry_in":6},"general":{"dir_paths":[],"heartbeat_interval":30},"db":{"type":1,"connection_string":"mongodb://localhost:27017"},"ffmpeg":{"use_double_quotes_for_path":false,"max_operation_retry_count":10000000,"rtmp_server_init_interval":3,"watch_dog_interval":23,"watch_dog_failed_wait_interval":3,"start_task_wait_for_interval":1,"record_concat_limit":1,"record_video_file_indexer_interval":60,"rtmp_server_port_start":7000,"rtmp_server_port_end":8000},"ai":{"overlay":true,"video_clip_duration":10,"face_recog_mtcnn_threshold":0.86,"face_recog_prob_threshold":0.98,"plate_recog_instance_count":2},"ui":{"gs_width":4,"gs_height":2,"booster_interval":0.3,"seek_to_live_edge_internal":30},"jobs":{"mac_ip_matching_enabled":false,"mac_ip_matching_interval":120,"black_screen_monitor_enabled":false,"black_screen_monitor_interval":600},"deep_stack":{"server_url":"http://127.0.0.1","server_port":1009,"performance_mode":1,"api_key":"","od_enabled":true,"od_threshold":0.45,"fr_enabled":true,"fr_threshold":0.7,"docker_type":1},"archive":{"limit_percent":95,"action_type":0,"move_location":""},"snapshot":{"process_count":1}}`

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
