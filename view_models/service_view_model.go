package view_models

import (
	"mngr/models"
	"mngr/utils"
	"time"
)

type ServiceViewModel struct {
	*models.ServiceModel
	RestartButtonEnabled bool `json:"restart_button_enabled"`
	StartButtonEnabled   bool `json:"start_button_enabled"`
	StopButtonEnabled    bool `json:"stop_button_enabled"`
}

func (s *ServiceViewModel) SetupButtonEnabled(containers map[string]*models.DockerContainer) {
	serviceNotRunning := true
	if s.InstanceType == models.Systemd {
		now := utils.StringToTime(utils.TimeToString(time.Now(), true))
		hb := utils.StringToTime(s.Heartbeat)
		diff := now.Sub(hb)
		serviceNotRunning = diff.Seconds() > 60
	} else if s.InstanceType == models.Container && len(s.InstanceName) > 0 {
		container, ok := containers[s.InstanceName]
		if ok {
			serviceNotRunning = container.State != models.Running
		}
	}
	s.RestartButtonEnabled = !serviceNotRunning
	s.StartButtonEnabled = serviceNotRunning
	s.StopButtonEnabled = !serviceNotRunning
}
