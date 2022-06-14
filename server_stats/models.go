package server_stats

type CpuInfo struct {
	Count             int     `json:"cpu_count"`
	UserUsage         float64 `json:"user_usage"`
	UserUsageHuman    string  `json:"user_usage_human"`
	SystemUsage       float64 `json:"system_usage"`
	SystemUsageHuman  string  `json:"system_usage_human"`
	Idle              float64 `json:"idle"`
	IdleHuman         string  `json:"idle_human"`
	UsagePercent      float64 `json:"usage_percent"`
	UsagePercentHuman string  `json:"usage_percent_human"`
}

type MemoryInfo struct {
	Total             uint64  `json:"total"`
	TotalHuman        string  `json:"total_human"`
	Used              uint64  `json:"used"`
	UsedHuman         string  `json:"used_human"`
	Cached            uint64  `json:"cached"`
	CachedHuman       string  `json:"cached_human"`
	Free              uint64  `json:"free"`
	FreeHuman         string  `json:"free_human"`
	UsagePercent      float64 `json:"usage_percent"`
	UsagePercentHuman string  `json:"usage_percent_human"`
}

type DiskInfo struct {
	MountPoint        string  `json:"mount_point"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	TotalHuman        string  `json:"total_human"`
	Used              uint64  `json:"used"`
	UsedHuman         string  `json:"used_human"`
	Free              uint64  `json:"free"`
	FreeHuman         string  `json:"free_human"`
	UsagePercent      float64 `json:"usage_percent"`
	UsagePercentHuman string  `json:"usage_percent_human"`
}

type NetworkInfo struct {
	Download      uint64 `json:"download"`
	DownloadHuman string `json:"download_human"`
	Upload        uint64 `json:"upload"`
	UploadHuman   string `json:"upload_human"`
}
