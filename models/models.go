package models

type Source struct {
	TypeName    string `json:"type_name"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	RtspAddress string `json:"rtsp_address"`
	Id          string `json:"id"`
}

type VideoFile struct {
	SourceId   string `json:"source_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Size       int64  `json:"size"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}
