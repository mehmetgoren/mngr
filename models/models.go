package models

import "path"

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

func (v *VideoFile) SetPath() {
	v.Path = path.Join("/playback", v.SourceId, v.Year, v.Month, v.Day, v.Hour, v.Name)
}

func (v *VideoFile) GetAbsolutePath(root string) string {
	return path.Join(root, v.SourceId, v.Year, v.Month, v.Day, v.Hour, v.Name)
}

type FFmpegReaderModel struct {
	Name   string `json:"name"`
	Img    string `json:"img"`
	Source string `json:"source"`
}
