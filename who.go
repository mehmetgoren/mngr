package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/reps"
	"os"
	"strconv"
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

func FetchRtspTemplates(rb *reps.RepoBucket) {
	// todo: this operation will be fetched by the Hub Portal instead of redis local db
}

func ReadEnvVariables(rb *reps.RepoBucket) *models.GlobalModel {
	config, _ := rb.ConfigRep.GetConfig()
	if config == nil || len(config.General.RootFolderPath) == 0 {
		config, _ = rb.ConfigRep.RestoreConfig()
	}

	mongoDbCs := os.Getenv("MONGODB_CS")
	if len(mongoDbCs) > 0 {
		config.Db.ConnectionString = mongoDbCs
		log.Println("MONGODB_CS: " + mongoDbCs)
	} else {
		fmt.Println("MONGODB_CS not found")
	}

	deepStackDt := os.Getenv("DEEPSTACK_DT")
	if len(deepStackDt) > 0 {
		val, _ := strconv.Atoi(deepStackDt)
		config.DeepStack.DockerType = val
		fmt.Println("DEEPSTACK_DT: " + deepStackDt)
	} else {
		fmt.Println("DEEPSTACK_DT not found")
	}

	deepStackOd := os.Getenv("DEEPSTACK_OD")
	if len(deepStackOd) > 0 {
		config.DeepStack.OdEnabled = deepStackOd == "1"
		fmt.Println("DEEPSTACK_OD: " + deepStackOd)
	} else {
		fmt.Println("DEEPSTACK_OD not found")
	}

	deepStackFr := os.Getenv("DEEPSTACK_FR")
	if len(deepStackFr) > 0 {
		config.DeepStack.FrEnabled = deepStackFr == "1"
		fmt.Println("DEEPSTACK_FR: " + deepStackFr)
	} else {
		fmt.Println("DEEPSTACK_FR not found")
	}

	rootDirPath := os.Getenv("ROOT_DIR_PATH")
	if len(rootDirPath) > 0 {
		config.General.RootFolderPath = rootDirPath
		fmt.Println("ROOT_DIR_PATH: " + rootDirPath)
	} else {
		fmt.Println("ROOT_DIR_PATH not found")
	}

	snapShotProcessCount := os.Getenv("SNAPSHOT_PROC_COUNT")
	if len(snapShotProcessCount) > 0 {
		processCount, _ := strconv.Atoi(snapShotProcessCount)
		if processCount < 1 {
			processCount = 1
		}
		config.Snapshot.ProcessCount = processCount
		fmt.Println("SNAPSHOT_PROC_COUNT: " + snapShotProcessCount)
	} else {
		fmt.Println("SNAPSHOT_PROC_COUNT not found")
	}

	rtmpServerPortStartStr := os.Getenv("RTMP_PORT_START")
	if len(rtmpServerPortStartStr) > 0 {
		rtmpServerPortStart, _ := strconv.Atoi(rtmpServerPortStartStr)
		if rtmpServerPortStart < 1 {
			rtmpServerPortStart = 1024
		}
		config.FFmpeg.RtmpServerPortStart = rtmpServerPortStart
		fmt.Println("RTMP_PORT_START: " + rtmpServerPortStartStr)
	} else {
		fmt.Println("RTMP_PORT_START not found")
	}

	rtmpServerPortEndStr := os.Getenv("RTMP_PORT_END")
	if len(rtmpServerPortEndStr) > 0 {
		rtmpServerPortEnd, _ := strconv.Atoi(rtmpServerPortEndStr)
		if rtmpServerPortEnd <= config.FFmpeg.RtmpServerPortStart {
			rtmpServerPortEnd = 65535
		}
		config.FFmpeg.RtmpServerPortEnd = rtmpServerPortEnd
		fmt.Println("RTMP_PORT_END: " + rtmpServerPortEndStr)
	} else {
		fmt.Println("RTMP_PORT_END not found")
	}

	global := &models.GlobalModel{}
	readOnlyModeStr := os.Getenv("READONLY_MODE")
	if len(readOnlyModeStr) > 0 {
		global.ReadOnlyMode = readOnlyModeStr == "1"
		fmt.Println("READONLY_MODE: " + readOnlyModeStr)
	} else {
		fmt.Println("READONLY_MODE not found")
	}

	rb.ConfigRep.SaveConfig(config)

	return global
}
