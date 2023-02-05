package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

func ReadEnvVariables(rb *reps.RepoBucket) *models.GlobalModel {
	config, _ := rb.ConfigRep.GetConfig()
	if config == nil || config.General.HeartbeatInterval == 0 { //config.General.HeartbeatInterval == 0  means that the config is not initialized
		config, _ = rb.ConfigRep.RestoreConfig()
	}

	rootDirPath := os.Getenv("ROOT_DIR_PATHS")
	if len(rootDirPath) > 0 {
		splits := strings.Split(rootDirPath, ",")
		l := len(splits)
		for i := 0; i < l; i++ {
			splits[i] = strings.TrimSpace(splits[i])
		}
		config.General.DirPaths = splits
		fmt.Println("ROOT_DIR_PATHS: " + rootDirPath)
	} else {
		fmt.Println("ROOT_DIR_PATHS not found")
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

func CheckSourceDirPaths(config *models.Config, rb *reps.RepoBucket) {
	setRootDir := func(c *models.Config, sourceId string) {
		sourceDirPath := utils.GetDefaultDirPath(c)
		s, _ := rb.SourceRep.Get(sourceId)
		s.RootDirPath = sourceDirPath
		rb.SourceRep.Save(s)
	}
	sources, _ := rb.SourceRep.GetAll()
	for _, source := range sources {
		if len(source.RootDirPath) == 0 {
			setRootDir(config, source.Id)
		} else if _, err := os.Stat(source.RootDirPath); os.IsNotExist(err) {
			setRootDir(config, source.Id)
		}
	}
}

func CheckHubContinuous(config *models.Config, rb *reps.RepoBucket, port int) {
	if !config.Hub.Enabled || len(config.Hub.Address) == 0 && len(config.Hub.Token) == 0 {
		log.Println("Hub integration is disabled")
		return
	}
	log.Println("Hub integration is enabled")
	maxRetry := config.Hub.MaxRetry
	if maxRetry < 1 {
		maxRetry = 100
	}
	count := 0
	for {
		if checkHub(config, rb, port) {
			if fetchRtspTemplates(config, rb) {
				log.Println("RTSP templates fetched from hub")
				break
			} else {
				log.Println("Failed to fetch RTSP templates from hub, retrying in 1 minute")
			}
		}
		if count > config.Hub.MaxRetry {
			log.Println("Hub max retry count reached, aborting")
			break
		}
		time.Sleep(time.Minute)
		count++
		log.Println("Hub integration retry count: " + strconv.Itoa(count))
	}
	log.Println("Hub integration completed")
}

func checkHub(config *models.Config, rb *reps.RepoBucket, port int) bool {
	ip, err := utils.GetExternalIP()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if len(ip) == 0 {
		log.Println("External IP not found")
		return false
	}
	var adminUser *models.User
	users, err := rb.UserRep.GetUsers()
	if err != nil {
		log.Println("an error occurred while gathering users", err.Error())
		return false
	}
	for _, user := range users {
		if user.Username == "admin" {
			adminUser = user
			break
		}
	}
	if adminUser == nil {
		log.Println("admin user not found")
		return false
	}

	requestModel := &models.NodeActivationRequest{
		NodeAddress:   "http://" + ip + ":" + strconv.Itoa(port),
		HubToken:      config.Hub.Token,
		NodeToken:     adminUser.Token,
		WebAppAddress: config.Hub.WebAppAddress,
	}

	url := config.Hub.Address + "/api/v1/node/activate"
	jsonBytes, err := json.Marshal(requestModel)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("an error an error occurred while creating a request", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("an error occurred while sending request", err.Error())
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(resp.Body)
	response := &models.NodeActivationResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		log.Println("an error occurred while decoding response", err.Error())
		return false
	}

	if response.Success {
		rb.AddUser(adminUser)
	}

	return response.Success
}

func fetchRtspTemplates(config *models.Config, rb *reps.RepoBucket) bool {
	url := config.Hub.Address + "/api/v1/node/getrtsptemplates"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("an error an error occurred while creating a request", err.Error())
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("an error occurred while sending request", err.Error())
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(resp.Body)
	response := make([]*models.RtspTemplateModel, 0)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("an error occurred while decoding response", err.Error())
		return false
	}

	_, err = rb.RtspTemplateRep.SaveAll(response)
	if err != nil {
		log.Println("an error occurred while saving rtsp templates", err.Error())
	}

	return true
}
