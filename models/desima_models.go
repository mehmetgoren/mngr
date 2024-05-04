package models

type NodeActivationRequest struct {
	DesimaToken   string `json:"DesimaToken"`
	NodeAddress   string `json:"NodeAddress"`
	NodeToken     string `json:"NodeToken"`
	WebAppAddress string `json:"WebAppAddress"`
}

type NodeActivationResponse struct {
	Success bool `json:"Success"`
}

type NodeDto struct {
	Id              string `json:"Id"`
	DesimaToken     string `json:"DesimaToken"`
	NodeAddress     string `json:"NodeAddress"`
	NodeToken       string `json:"NodeToken"`
	WebAppAddress   string `json:"WebAppAddress"`
	Name            string `json:"Name"`
	Enabled         bool   `json:"Enabled"`
	Description     string `json:"Description"`
	IsActivated     bool   `json:"IsActivated"`
	Available       bool   `json:"Available"`
	UptimeInSeconds int64  `json:"UptimeInSeconds"`
	OperatingSystem string `json:"OperatingSystem"`
	SourceCount     int    `json:"SourceCount"`
	StreamCount     int    `json:"StreamCount"`

	CpuCount             int     `json:"CpuCount"`
	CpuUsagePercent      float64 `json:"CpuUsagePercent"`
	CpuUsagePercentHuman string  `json:"CpuUsagePercentHuman"`

	MemoryTotal             int64   `json:"MemoryTotal"`
	MemoryTotalHuman        string  `json:"MemoryTotalHuman"`
	MemoryUsed              int64   `json:"MemoryUsed"`
	MemoryUsedHuman         string  `json:"MemoryUsedHuman"`
	MemoryFree              int64   `json:"MemoryFree"`
	MemoryFreeHuman         string  `json:"MemoryFreeHuman"`
	MemoryUsagePercent      float64 `json:"MemoryUsagePercent"`
	MemoryUsagePercentHuman string  `json:"MemoryUsagePercentHuman"`

	DiskTotal             int64   `json:"DiskTotal"`
	DiskTotalHuman        string  `json:"DiskTotalHuman"`
	DiskUsed              int64   `json:"DiskUsed"`
	DiskUsedHuman         string  `json:"DiskUsedHuman"`
	DiskFree              int64   `json:"DiskFree"`
	DiskFreeHuman         string  `json:"DiskFreeHuman"`
	DiskUsagePercent      float64 `json:"DiskUsagePercent"`
	DiskUsagePercentHuman string  `json:"DiskUsagePercentHuman"`

	GpuName          string `json:"GpuName"`
	GpuDriverVersion string `json:"GpuDriverVersion"`
	GpuCudaVersion   string `json:"GpuCudaVersion"`
	GpuMemoryTotal   string `json:"GpuMemoryTotal"`
	GpuMemoryUsed    string `json:"GpuMemoryUsed"`
	GpuPowerLimit    string `json:"GpuPowerLimit"`
	GpuPowerDraw     string `json:"GpuPowerDraw"`

	TotalAiDetection         int64  `json:"TotalAiDetection"`
	TotalRegisteredFaces     int    `json:"TotalRegisteredFaces"`
	RunningServices          string `json:"RunningServices"`
	GdriveEnabled            bool   `json:"GdriveEnabled"`
	TelegramEnabled          bool   `json:"TelegramEnabled"`
	TotalUserCount           int    `json:"TotalUserCount"`
	MsContainerFailedCount   int    `json:"MsContainerFailedCount"`
	MsFeederFailedCount      int    `json:"MsFeederFailedCount"`
	HlsFailedCount           int    `json:"HlsFailedCount"`
	FfmpegReaderFailedCount  int    `json:"FfmpegReaderFailedCount"`
	RecordFailedCount        int    `json:"RecordFailedCount"`
	SnapshotFailedCount      int    `json:"SnapshotFailedCount"`
	RecordStuckProcessCount  int    `json:"RecordStuckProcessCount"`
	SourceStateConflictCount int    `json:"SourceStateConflictCount"`
}
