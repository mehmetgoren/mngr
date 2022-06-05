package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/reps"
)

func createHeartbeatRepository(client *redis.Client, serviceName string, config *models.Config) *reps.HeartbeatRepository {
	var heartbeatRepository = reps.HeartbeatRepository{Connection: client, TimeSecond: int64(config.General.HeartbeatInterval), ServiceName: serviceName}

	return &heartbeatRepository
}

func WhoAreYou(rb *reps.RepoBucket) {
	connMain := rb.GetMainConnection()
	config, _ := rb.ConfigRep.GetConfig()

	serviceName := "web_server_manager"
	heartbeat := createHeartbeatRepository(connMain, serviceName, config)
	go heartbeat.Start()

	serviceRepository := rb.ServiceRep
	go func() {
		_, err := serviceRepository.Add(serviceName)
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()
}
