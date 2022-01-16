package models

type Source struct {
	TypeName    string `json:"type_name"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	RtspAddress string `json:"rtsp_address"`
	Id          string `json:"id"`
}

type Recording struct {
	Id          string `json:"id" redis:"id"`
	Name        string `json:"name" redis:"name"`
	Brand       string `json:"brand" redis:"brand"`
	RtspAddress string `json:"rtsp_address" redis:"rtsp_address"`
	Pid         string `json:"pid" redis:"pid"`
	OutputFile  string `json:"output_file" redis:"output_file"`
	FailedCount int    `json:"failed_count" redis:"failed_count"`
	CreatedAt   string `json:"created_at" redis:"created_at"`
	Duration    int    `json:"duration" redis:"duration"`
	Args        string `json:"args" redis:"args"`
}

type VideoFile struct {
	SourceId   string  `json:"source_id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Size       float64 `json:"size"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
}
