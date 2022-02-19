package main

import (
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
)

const (
	MAIN     = 0
	RQ       = 1
	EVENTBUS = 15
)

func createHeartbeatRepository(client *redis.Client, serviceName string) *reps.HeartbeatRepository {
	var heartbeatRepository = reps.HeartbeatRepository{Client: client, TimeSecond: 10, ServiceName: serviceName}

	return &heartbeatRepository
}

func createServiceRepository(client *redis.Client) *reps.ServiceRepository {
	var pidRepository = reps.ServiceRepository{Client: client}

	return &pidRepository
}

func WhoAreYou() {
	client := reps.CreateRedisConnection(MAIN)

	serviceName := "web_server_manager"
	heartbeat := createHeartbeatRepository(client, serviceName)
	go heartbeat.Start()

	serviceRepository := createServiceRepository(client)
	go func() {
		_, err := serviceRepository.Add(serviceName)
		if err != nil {
			log.Println("An error occurred while registering process id, error is:" + err.Error())
		}
	}()
}
