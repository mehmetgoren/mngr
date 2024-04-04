package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func WhoAreYou(rb *reps.RepoBucket) {
	serviceName := "web_server_manager"
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
	if config == nil || config.Archive.LimitPercent == 0 { //config.Archive.LimitPercent == 0  means that the config is not initialized
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

	// todo: Move this env variable to the wizard
	senseAiImage := os.Getenv("SENSE_AI_IMAGE")
	if len(senseAiImage) > 0 {
		img, _ := strconv.Atoi(senseAiImage)
		config.SenseAi.Image = img
		fmt.Println("SENSE_AI_IMAGE: " + strconv.Itoa(img))
	} else {
		fmt.Println("SENSE_AI_IMAGE not found")
	}

	senseAiHost := os.Getenv("SENSE_AI_HOST")
	if len(senseAiHost) > 0 {
		config.SenseAi.Host = senseAiHost
		fmt.Println("SENSE_AI_HOST: " + senseAiHost)
	} else {
		fmt.Println("SENSE_AI_HOST not found")

	}
	senseAiPort := os.Getenv("SENSE_AI_PORT")
	if len(senseAiPort) > 0 {
		port, _ := strconv.Atoi(senseAiPort)
		if port > 0 {
			config.SenseAi.Port = port
		}
		fmt.Println("SENSE_AI_PORT: " + strconv.Itoa(port))
	} else {
		fmt.Println("SENSE_AI_PORT not found")

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

	msPortStartStr := os.Getenv("MS_PORT_START")
	if len(msPortStartStr) > 0 {
		msPortStart, _ := strconv.Atoi(msPortStartStr)
		if msPortStart < 1 {
			msPortStart = 1024
		}
		config.FFmpeg.MsPortStart = msPortStart
		fmt.Println("MS_PORT_START: " + msPortStartStr)
	} else {
		fmt.Println("MS_PORT_START not found")
	}

	msPortEndStr := os.Getenv("MS_PORT_END")
	if len(msPortEndStr) > 0 {
		msPortEnd, _ := strconv.Atoi(msPortEndStr)
		if msPortEnd <= config.FFmpeg.MsPortStart {
			msPortEnd = 65535
		}
		config.FFmpeg.MsPortEnd = msPortEnd
		fmt.Println("MS_PORT_END: " + msPortEndStr)
	} else {
		fmt.Println("MS_PORT_END not found")
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
	dirPathsMap := make(map[string]bool)
	for _, dirPath := range config.General.DirPaths {
		dirPathsMap[dirPath] = true
	}

	sources, _ := rb.SourceRep.GetAll()
	for _, source := range sources {
		if len(source.RootDirPath) == 0 {
			setRootDir(config, source.Id)
			log.Println("RootDirPath is not set for source: "+source.Id, source.Name)
		} else if !utils.IsDirExists(source.RootDirPath) {
			setRootDir(config, source.Id)
			log.Println("RootDirPath is not exists for source: "+source.Id, source.Name)
		} else if _, ok := dirPathsMap[source.RootDirPath]; !ok {
			setRootDir(config, source.Id)
			log.Println("RootDirPath is not exists in Config.General.DirPaths for source: "+source.Id, source.Name)
		} else {
			log.Println("RootDirPath is ok for: "+source.Id, source.Name)
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

func CheckMissingConfigValues(config *models.Config, rb *reps.RepoBucket) {
	if config.General.DirPaths == nil || len(config.General.DirPaths) == 0 {
		log.Println("No directory path is set, the program will be terminated")
		os.Exit(1)
		return
	}
	for _, dirPath := range config.General.DirPaths {
		if !utils.IsDirExists(dirPath) {
			log.Println("The directory path (" + dirPath + ") does not exist, the program will be terminated")
			os.Exit(1)
			return
		}
	}
	log.Println("DirPaths: ", config.General.DirPaths)

	originalConfig, _ := utils.DeepCopy(config)

	if len(config.Db.ConnectionString) == 0 {
		config.Db.Type = 1 // 0:SQLite, 1: mongodb
		config.Db.ConnectionString = "mongodb://localhost:27017"
	}

	if config.FFmpeg.MaxOperationRetryCount == 0 {
		config.FFmpeg.MaxOperationRetryCount = 10000000
	}
	if config.FFmpeg.MsInitInterval == 0 {
		config.FFmpeg.MsInitInterval = 3
	}
	if config.FFmpeg.WatchDogInterval == 0 {
		config.FFmpeg.WatchDogInterval = 23
	}
	if config.FFmpeg.WatchDogFailedWaitInterval == 0 {
		config.FFmpeg.WatchDogFailedWaitInterval = 3
	}
	if config.FFmpeg.StartTaskWaitForInterval == 0 {
		config.FFmpeg.StartTaskWaitForInterval = 1
	}
	if config.FFmpeg.RecordConcatLimit == 0 {
		config.FFmpeg.RecordConcatLimit = 1
	}
	if config.FFmpeg.RecordVideoFileIndexerInterval == 0 {
		config.FFmpeg.RecordVideoFileIndexerInterval = 60
	}
	if config.FFmpeg.MsPortStart == 0 {
		config.FFmpeg.MsPortStart = 7000
	}
	if config.FFmpeg.MsPortEnd == 0 {
		config.FFmpeg.MsPortEnd = 8000
	}

	if config.Ai.VideoClipDuration == 0 {
		config.Ai.VideoClipDuration = 10
	}

	if len(config.SenseAi.Host) == 0 {
		config.SenseAi.Host = "127.0.0.1"
	}
	if config.SenseAi.Port == 0 {
		config.SenseAi.Port = 32168
	}

	if config.Jobs.MacIpMatchingInterval == 0 {
		config.Jobs.MacIpMatchingInterval = 120
	}
	if config.Jobs.BlackScreenMonitorInterval == 0 {
		config.Jobs.BlackScreenMonitorInterval = 600
	}

	if config.Archive.LimitPercent == 0 {
		config.Archive.LimitPercent = 95
	}

	if config.Snapshot.ProcessCount == 0 {
		config.Snapshot.ProcessCount = 4
	}

	if len(config.Hub.Address) == 0 {
		config.Hub.Address = "http://localhost:5268"
	}
	if len(config.Hub.WebAppAddress) == 0 {
		config.Hub.WebAppAddress = "http://localhost:8080"
	}
	if config.Hub.MaxRetry == 0 {
		config.Hub.MaxRetry = 100
	}

	if !reflect.DeepEqual(originalConfig, config) {
		err := rb.ConfigRep.SaveConfig(config)
		if err != nil {
			log.Println("Error while saving config file: ", err)
		}
		log.Println("Config file is updated")
	} else {
		log.Println("Config file is not changed")
	}
}
