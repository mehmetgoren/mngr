package models

const (
	Created    = "created"
	Running    = "running"
	Paused     = "paused"
	Restarting = "restarting"
	Removing   = "removing"
	Exited     = "exited"
	Dead       = "dead"
)

type DockerContainer struct {
	Name  string `json:"name"`
	State string `json:"state"`
}
