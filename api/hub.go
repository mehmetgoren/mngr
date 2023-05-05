package api

import (
	human "github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"
	"log"
	"mngr/data"
	"mngr/data/cmn"
	"mngr/models"
	"mngr/reps"
	"mngr/server_stats"
	"mngr/utils"
	"mngr/ws"
	"net/http"
	"os"
	"strings"
	"time"
)

func RegisterHubEndpoints(router *gin.Engine, rb *reps.RepoBucket, holders *ws.Holders, factory *cmn.Factory) {
	router.POST("/loginbytoken", func(ctx *gin.Context) {
		var lu models.LoginUserByTokenViewModel
		if err := ctx.BindJSON(&lu); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := rb.UserRep.LoginByToken(lu.Token)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if u == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "401 Unauthorized"})
			return
		}
		logoutUser := func(user *models.User, triggerLogout bool) {
			rb.RemoveUser(user.Token)
			holders.UserLogout(user.Token, triggerLogout)
		}

		logoutUser(u, false)
		time.Sleep(1 * time.Second)
		if u != nil {
			rb.AddUser(u)
		}
		ctx.JSON(http.StatusOK, u)
	})

	router.POST("/nodeinfo", func(ctx *gin.Context) {
		var p models.NodeDto
		if err := ctx.ShouldBindJSON(&p); err != nil {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		nodeDto := &p

		nodeDto.UptimeInSeconds = int64(time.Now().Sub(utils.StartupTime).Seconds())
		sources, _ := rb.SourceRep.GetAll()
		if sources != nil {
			nodeDto.SourceCount = len(sources)
		}
		serviceModel, _ := rb.ServiceRep.GetServices()
		var me *models.ServiceModel
		if serviceModel != nil && len(serviceModel) > 0 {
			for _, service := range serviceModel {
				if service.Name == "web_application" {
					me = service
					break
				}
			}
		}
		if me != nil {
			nodeDto.OperatingSystem = me.Platform
		}
		config, _ := rb.ConfigRep.GetConfig()
		nodeDto.WebAppAddress = config.Hub.WebAppAddress
		setHardwareInfo(nodeDto, config)
		setGpuInfo(nodeDto)
		setAiInfo(nodeDto, config, factory)
		setServicesInfo(nodeDto, rb)
		setCloudInfo(nodeDto, rb)
		setUsersInfo(nodeDto, rb)
		setFailedInfo(nodeDto, rb)

		ctx.JSON(http.StatusOK, nodeDto)
	})
}

func setHardwareInfo(nodeDto *models.NodeDto, config *models.Config) {
	stats := &server_stats.ServerStats{}
	err := stats.InitCpuInfos()
	if err == nil {
		nodeDto.CpuCount = stats.Cpu.Count
		nodeDto.CpuUsagePercent = stats.Cpu.UsagePercent
		nodeDto.CpuUsagePercentHuman = stats.Cpu.UsagePercentHuman
	}

	err = stats.InitMemInfos()
	if err == nil {
		nodeDto.MemoryTotal = int64(stats.Memory.Total)
		nodeDto.MemoryTotalHuman = stats.Memory.TotalHuman
		nodeDto.MemoryUsed = int64(stats.Memory.Used)
		nodeDto.MemoryUsedHuman = stats.Memory.UsedHuman
		nodeDto.MemoryFree = int64(stats.Memory.Free)
		nodeDto.MemoryFreeHuman = stats.Memory.FreeHuman
		nodeDto.MemoryUsagePercent = stats.Memory.UsagePercent
		nodeDto.MemoryUsagePercentHuman = stats.Memory.UsagePercentHuman
	}

	err = stats.InitDiskInfos(config)
	diskCount := len(stats.Disks)
	if err == nil && stats.Disks != nil && diskCount > 0 {
		diskTotal := uint64(0)
		diskUsed := uint64(0)
		diskFree := uint64(0)
		for _, disk := range stats.Disks {
			diskTotal += disk.Total
			diskUsed += disk.Used
			diskFree += disk.Free
		}
		nodeDto.DiskTotal = int64(diskTotal)
		nodeDto.DiskTotalHuman = human.Bytes(diskTotal * 1024 * 1024)
		nodeDto.DiskUsed = int64(diskUsed)
		nodeDto.DiskUsedHuman = human.Bytes(diskUsed * 1024 * 1024)
		nodeDto.DiskFree = int64(diskFree)
		nodeDto.DiskFreeHuman = human.Bytes(diskFree)
		nodeDto.DiskUsagePercent = (float64(diskUsed) / float64(diskTotal)) * 100.0
		nodeDto.DiskUsagePercentHuman = human.CommafWithDigits(nodeDto.DiskUsagePercent, 2) + " %"
	}
}

func setGpuInfo(nodeDto *models.NodeDto) {
	gpuModel := utils.NvidiaGpuModel{}
	err := gpuModel.Fetch()
	if err == nil {
		nodeDto.GpuName = gpuModel.ProductName
		nodeDto.GpuDriverVersion = gpuModel.DriverVersion
		nodeDto.GpuCudaVersion = gpuModel.CudaVersion
		nodeDto.GpuMemoryTotal = gpuModel.MemoryTotal
		nodeDto.GpuMemoryUsed = gpuModel.MemoryUsed
		nodeDto.GpuPowerLimit = gpuModel.PowerLimit
		nodeDto.GpuPowerDraw = gpuModel.PowerDraw
	}
}

func setAiInfo(nodeDto *models.NodeDto, config *models.Config, factory *cmn.Factory) {
	startDate := time.Date(2022, 10, 25, 0, 0, 0, 0, time.UTC)
	endDate := time.Now()

	queryParams := data.QueryParams{}
	queryParams.T1 = startDate
	queryParams.T2 = endDate

	rep := factory.CreateRepository()

	nodeDto.TotalObjectDetection, _ = rep.CountOds(queryParams)
	nodeDto.TotalFaceDetection, _ = rep.CountFrs(queryParams)
	nodeDto.TotalLicensePlateDetection, _ = rep.CountAlprs(queryParams)

	mlTrainPath := utils.GetFrTrainPath(config)
	de, err := os.ReadDir(mlTrainPath)
	if err == nil {
		totalRegisteredFaces := 0
		for _, d := range de {
			if d.IsDir() {
				totalRegisteredFaces++
			}
		}
		nodeDto.TotalRegisteredFaces = totalRegisteredFaces
	}
}

func setServicesInfo(nodeDto *models.NodeDto, rb *reps.RepoBucket) {
	services, err := rb.ServiceRep.GetServices()
	if err != nil {
		log.Println("Error while getting services, err: ", err.Error())
		return
	}
	var sb strings.Builder
	l := len(services)
	j := 0
	for ; j < l-1; j++ {
		sb.WriteString(services[j].Description)
		sb.WriteString(",")
	}
	sb.WriteString(services[j].Description)
	nodeDto.RunningServices = sb.String()
}

func setCloudInfo(nodeDto *models.NodeDto, rb *reps.RepoBucket) {
	rep := rb.CloudRep
	nodeDto.TelegramEnabled = rep.IsTelegramIntegrationEnabled()
	nodeDto.GdriveEnabled = rep.IsGdriveIntegrationEnabled()
}

func setUsersInfo(nodeDto *models.NodeDto, rb *reps.RepoBucket) {
	user, _ := rb.UserRep.GetUsers()
	if user != nil {
		users, _ := rb.UserRep.GetUsers()
		if users != nil {
			nodeDto.TotalUserCount = len(users)
		}
	}
}

func setFailedInfo(nodeDto *models.NodeDto, rb *reps.RepoBucket) {
	rep := rb.FailedStreamRep
	failedStreams, _ := rep.GetAll()
	if failedStreams != nil {
		for _, failedStream := range failedStreams {
			nodeDto.MsContainerFailedCount += failedStream.MsContainerFailedCount
			nodeDto.MsFeederFailedCount += failedStream.MsFeederFailedCount
			nodeDto.HlsFailedCount += failedStream.HlsFailedCount
			nodeDto.FfmpegReaderFailedCount += failedStream.FFmpegReaderFailedCount
			nodeDto.RecordFailedCount += failedStream.RecordFailedCount
			nodeDto.SnapshotFailedCount += failedStream.SnapshotFailedCount
			nodeDto.RecordStuckProcessCount += failedStream.RecordStuckProcessCount
			nodeDto.SourceStateConflictCount += failedStream.SourceStateConflictCount
		}
	}
}
