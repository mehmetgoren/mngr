package models

type OnvifParams struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Asset struct {
	Address string `json:"address"`
}

type OnvifEvent struct {
	Type        string `json:"type"`
	Base64Model string `json:"base64_model"`
}

type ExecResultView struct {
	Device   string   `json:"device"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Routes   []string `json:"route"`
	Address  string   `json:"address"`
	Port     uint16   `json:"port"`

	CredentialsFound bool `json:"credentials_found"`
	RouteFound       bool `json:"route_found"`
	Available        bool `json:"available"`

	AuthenticationType string `json:"authentication_type"`
}

type NetworkDiscoveryModel struct {
	Results   []*ExecResultView `json:"results"`
	CreatedAt string            `json:"created_at"`
}

type DeviceInfo struct {
	Manufacturer    string `json:"manufacturer"`
	Model           string `json:"model"`
	FirmwareVersion string `json:"firmware_version"`
	SerialNumber    string `json:"serial_number"`
	HardwareId      string `json:"hardware_id"`
}

type User struct {
	Username  string `json:"username"`
	UserLevel string `json:"user_level"`
}

type TargetInfo struct {
	DeviceInfo *DeviceInfo `json:"device_info"`

	IsDiscoverable bool   `json:"is_discoverable"`
	HostName       string `json:"host_name"`

	IPAddresses []string `json:"ip_addresses"`

	//GetNetworkInterfaces
	HwAddress string `json:"hw_address"`

	//GetNetworkProtocols
	HttpPort  int `json:"http_port"`
	HttpsPort int `json:"https_port"`
	RtspPort  int `json:"rtsp_port"`

	//GetSystemDateAndTime
	LocalDatetime string `json:"local_datetime"`

	//GetSystemLog
	Logs string `json:"logs"`

	Users []*User `json:"users"`

	StreamUri string `json:"stream_uri"`
}

type OnvifModel struct {
	HackResult  *ExecResultView `json:"hack_result"`
	Onvif       *TargetInfo     `json:"onvif"`
	OnvifParams *OnvifParams    `json:"onvif_params"`
	CreatedAt   string          `json:"created_at"`
}
