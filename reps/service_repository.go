package reps

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
	"os"
	"runtime"
	"time"
)

type ServiceRepository struct {
	Connection *redis.Client
}

func getServicesKey(serviceName string) string {
	return "services:" + serviceName
}

func (r *ServiceRepository) Add(serviceName string) (int64, error) {
	now := time.Now()
	host, _ := os.Hostname()
	sm := &models.ServiceModel{
		Name:            serviceName,
		Description:     "The Web Server Manager Service®",
		Platform:        runtime.GOOS,
		PlatformVersion: runtime.GOARCH,
		HostName:        host,
		IpAddress:       "127.0.0.1",
		MacAddress:      "00:00:00:00:00:00",
		Processor:       "unknown",
		CpuCount:        runtime.NumCPU(),
		Ram:             "unknown",
		Pid:             os.Getpid(),
		CreatedAt:       utils.TimeToString(now, true),
		Heartbeat:       "",
	}

	return r.Connection.HSet(context.Background(), getServicesKey(serviceName), Map(sm)).Result()
}

func (r *ServiceRepository) GetServices() ([]*models.ServiceModel, error) {
	conn := r.Connection
	list := make([]*models.ServiceModel, 0, 5)
	keys, err := conn.Keys(context.Background(), getServicesKey("*")).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			conn.Set(context.Background(), redisKeyStream, list, 0)
			return list, nil
		} else {
			log.Println("Error getting all stream from redis: ", err)
			return nil, err
		}
	}

	for _, key := range keys {
		var p models.ServiceModel
		err := conn.HGetAll(context.Background(), key).Scan(&p)
		if err != nil {
			log.Println("Error getting stream from redis: ", err)
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}
