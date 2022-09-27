package models

type SourceStatusModel struct {
	SourceId  string `json:"id"`
	Enabled   bool   `json:"enabled"`
	Streaming bool   `json:"streaming"`
	Recording bool   `json:"recording"`
}

type SourceEnabledModel struct {
	SourceId string `json:"id"`
	Enabled  bool   `json:"enabled"`
}
