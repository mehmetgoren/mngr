package models

type VideoFile struct {
	SourceId   string  `json:"source_id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Size       float64 `json:"size"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
}

type FFmpegReaderModel struct {
	Name   string `json:"name"`
	Img    string `json:"img"`
	Source string `json:"source"`
}
