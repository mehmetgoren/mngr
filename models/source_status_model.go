package models

type SourceStatusModel struct {
	SourceId  string `json:"id"`
	Streaming bool   `json:"streaming"`
}
