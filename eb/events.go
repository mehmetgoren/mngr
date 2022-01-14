package eb

import (
	"encoding/json"
	"mngr/models"
)

type Event struct {
}

type StreamingEvent struct {
	models.Source
	OutputFile string `json:"output_file"`
}

func (s StreamingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s StreamingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

type RecordingEvent struct {
	models.Source
	Duration   int    `json:"duration"`
	OutputFile string `json:"output_file"`
}

func (s RecordingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s RecordingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
