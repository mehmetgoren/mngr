package models

type FFmpegModel interface {
	GetSourceId() string
	GetDirPath() string
}
