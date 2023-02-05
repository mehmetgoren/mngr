package models

type InstanceType int

const (
	Systemd   InstanceType = 0
	Container InstanceType = 1
)

type ServiceModel struct {
	Name            string       `json:"name" redis:"name"`
	Description     string       `json:"description" redis:"description"`
	Platform        string       `json:"platform" redis:"platform"`
	PlatformVersion string       `json:"platform_version" redis:"platform_version"`
	HostName        string       `json:"hostname" redis:"hostname"`
	IpAddress       string       `json:"ip_address" redis:"ip_address"`
	MacAddress      string       `json:"mac_address" redis:"mac_address"`
	Processor       string       `json:"processor" redis:"processor"`
	CpuCount        int          `json:"cpu_count" redis:"cpu_count"`
	Ram             string       `json:"ram" redis:"ram"`
	Pid             int          `json:"pid" redis:"pid"`
	InstanceType    InstanceType `json:"instance_type" redis:"instance_type"`
	InstanceName    string       `json:"instance_name" redis:"instance_name"`
	CreatedAt       string       `json:"created_at" redis:"created_at"`
	Heartbeat       string       `json:"heartbeat" redis:"heartbeat"`
}

type RegisterWebAppServiceModel struct {
	AppAddress string `json:"app_address"`
}
