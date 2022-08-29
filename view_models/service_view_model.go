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

func (s *ServiceViewModel) SetupButtonEnabled() {
	now := utils.StringToTime(utils.TimeToString(time.Now(), true))
	hb := utils.StringToTime(s.Heartbeat)
	diff := now.Sub(hb)
	serviceNotRunning := diff.Seconds() > 60
	s.RestartButtonEnabled = !serviceNotRunning
	s.StartButtonEnabled = serviceNotRunning
	s.StopButtonEnabled = !serviceNotRunning
}
