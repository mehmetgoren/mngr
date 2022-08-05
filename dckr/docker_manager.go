package dckr

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
	"mngr/models"
)

type DockerManager struct {
	Client *client.Client
}

func (d *DockerManager) getContainer(instanceName string) (*types.Container, error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{All: false})
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

func (d *DockerManager) RestartAfterCloudChanges() bool {
	if d.Client == nil {
		return false
	}
	names := []string{"smcp-instance"}
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

func (d *DockerManager) RestartAll(services []*models.ServiceModel) bool {
	if d.Client == nil || services == nil {
		return false
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
	}
	return true
}
