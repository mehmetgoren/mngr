package api

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"log"
	"mngr/dckr"
	"mngr/models"
	"mngr/reps"
	"mngr/view_models"
	"net/http"
	"strings"
)

func RegisterServiceEndpoints(router *gin.Engine, rb *reps.RepoBucket, dockerClient *client.Client) {
	router.GET("/services", func(ctx *gin.Context) {
		services, err := rb.ServiceRep.GetServices()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		vms := make([]*view_models.ServiceViewModel, 0)
		dm := dckr.DockerManager{Client: dockerClient}
		containers, _ := dm.GetContainers()
		for _, service := range services {
			vm := &view_models.ServiceViewModel{ServiceModel: service}
			vm.SetupButtonEnabled(containers)
			vms = append(vms, vm)
		}
		ctx.JSON(http.StatusOK, vms)
	})

	router.POST("/registerwebappservice", func(ctx *gin.Context) {
		var p models.RegisterWebAppServiceModel
		if err := ctx.ShouldBindJSON(&p); err != nil {
			log.Println("Got error while unmarshalling RegisterWebAppServiceModel, err", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(p.AppAddress) > 0 {
			c, _ := rb.ConfigRep.GetConfig()
			if len(c.Hub.WebAppAddress) == 0 {
				address := p.AppAddress
				address = strings.ReplaceAll(address, "#/", "")
				splits := strings.Split(address, "?")
				if len(splits) > 0 {
					address = splits[0]
					if strings.HasSuffix(address, "/") {
						address = address[:len(address)-1]
					}
				}
				c.Hub.WebAppAddress = address
				err := rb.ConfigRep.SaveConfig(c)
				if err != nil {
					log.Println("Error while saving config: ", err.Error())
				}
			}
		}

		_, err := rb.ServiceRep.AddWebApp("web_application")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, true)
	})

	router.POST("/restartservice", func(ctx *gin.Context) {
		containerAction(dockerClient, ctx, func(service *models.ServiceModel, dm *dckr.DockerManager) bool {
			return dm.RestartContainer(service.InstanceName)
		})
	})

	router.POST("/startservice", func(ctx *gin.Context) {
		containerAction(dockerClient, ctx, func(service *models.ServiceModel, dm *dckr.DockerManager) bool {
			return dm.StartContainer(service.InstanceName)
		})
	})

	router.POST("/stopservice", func(ctx *gin.Context) {
		containerAction(dockerClient, ctx, func(service *models.ServiceModel, dm *dckr.DockerManager) bool {
			return dm.StopContainer(service.InstanceName)
		})
	})

	router.POST("/restartaftercloudchanges", func(ctx *gin.Context) {
		dm := dckr.DockerManager{Client: dockerClient}
		result := dm.RestartAfterCloudChanges()
		ctx.JSON(http.StatusOK, result)
	})

	router.POST("/restartallservices", func(ctx *gin.Context) {
		services, err := rb.ServiceRep.GetServices()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		dm := dckr.DockerManager{Client: dockerClient}
		result := dm.RestartAll(services)
		ctx.JSON(http.StatusOK, result)
	})
}

func containerAction(dockerClient *client.Client, ctx *gin.Context, fn func(sm *models.ServiceModel, dm *dckr.DockerManager) bool) {
	var service models.ServiceModel
	if err := ctx.ShouldBindJSON(&service); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if service.InstanceType == models.Systemd {
		ctx.JSON(http.StatusOK, false)
		return
	}

	dm := dckr.DockerManager{Client: dockerClient}
	result := fn(&service, &dm) // dm.RestartContainer(service.InstanceName)
	ctx.JSON(http.StatusOK, result)
}
