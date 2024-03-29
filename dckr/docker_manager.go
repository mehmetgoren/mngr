package dckr

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"mngr/models"
	"strings"
	"time"
)

type DockerManager struct {
	Client *client.Client
}

var smcpInstance = "smcp-instance"
var mngrInstance = "mngr-instance"

func (d *DockerManager) getContainer(instanceName string) (*types.Container, error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Println("an error occurred while getting the container, err: ", err.Error())
		return nil, err
	}

	instanceName = "/" + instanceName
	for _, cntr := range containers {
		for _, cname := range cntr.Names {
			if cname == instanceName {
				return &cntr, nil
			}
		}
	}

	return nil, nil
}

func (d *DockerManager) RestartContainer(instanceName string) bool {
	if d.Client == nil || len(instanceName) == 0 {
		return false
	}

	cntr, _ := d.getContainer(instanceName)
	if cntr == nil || cntr.State != "running" {
		return false
	}
	err := d.Client.ContainerRestart(context.Background(), cntr.ID, nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (d *DockerManager) StartContainer(instanceName string) bool {
	if d.Client == nil || len(instanceName) == 0 {
		return false
	}

	cntr, _ := d.getContainer(instanceName)
	if cntr == nil || cntr.State == "running" {
		return false
	}
	err := d.Client.ContainerStart(context.Background(), cntr.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (d *DockerManager) StopContainer(instanceName string) bool {
	if d.Client == nil || len(instanceName) == 0 {
		return false
	}

	cntr, _ := d.getContainer(instanceName)
	if cntr == nil || cntr.State != "running" {
		return false
	}
	err := d.Client.ContainerStop(context.Background(), cntr.ID, nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (d *DockerManager) RestartAfterCloudChanges() bool {
	if d.Client == nil {
		return false
	}
	names := []string{smcpInstance}
	for _, name := range names {
		cntr, _ := d.getContainer(name)
		if cntr == nil || cntr.State != "running" {
			continue
		}
		err := d.Client.ContainerRestart(context.Background(), cntr.ID, nil)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	return true
}

var isRestartAllRunning = false

func (d *DockerManager) RestartAll(services []*models.ServiceModel) bool {
	if isRestartAllRunning {
		return false
	}
	isRestartAllRunning = true
	defer func() {
		isRestartAllRunning = false
	}()
	if d.Client == nil || services == nil {
		return false
	}
	count := len(services)
	for index, service := range services {
		if service.InstanceName == mngrInstance {
			lastService := services[count-1]
			services[count-1] = service
			services[index] = lastService
			break
		}
	}
	for _, service := range services {
		if service.InstanceType == models.Systemd {
			continue
		}
		cntr, _ := d.getContainer(service.InstanceName)
		if cntr == nil || cntr.State != "running" {
			continue
		}
		err := d.Client.ContainerRestart(context.Background(), cntr.ID, nil)
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(time.Second)
	}
	return true
}

func (d *DockerManager) GetContainers() (map[string]*models.DockerContainer, error) {
	ret := make(map[string]*models.DockerContainer)
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Println("an error occurred while getting the container, err: ", err.Error())
		return ret, err
	}
	for _, c := range containers {
		if c.Names == nil || len(c.Names) == 0 {
			continue
		}
		name := c.Names[0]
		name = strings.Replace(name, "/", "", -1)
		dc := &models.DockerContainer{Name: name, State: c.State}
		ret[name] = dc
	}

	return ret, nil
}
