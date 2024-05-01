package models

type VideoFile struct {
	SourceId   string  `json:"source_id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Size       float64 `json:"size"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
	Year       string  `json:"year"`
	Month      string  `json:"month"`
	Day        string  `json:"day"`
	Hour       string  `json:"hour"`
}

type FFmpegReaderModel struct {
	Name        string `json:"name"`
	Base64Image string `json:"base64_image"`
	SourceId    string `json:"source_id"`
}
