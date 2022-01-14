package eb

import (
	"github.com/go-redis/redis/v8"
	"mngr/utils"
	"mngr/ws"
	"path/filepath"
)

func (s *StreamingEvent) Handle(event *redis.Message) error {
	var eventModel StreamingEvent
	utils.DeserializeJson(event.Payload, &eventModel)

	eventModel.OutputFile = filepath.Join(eventModel.Id, "stream.m3u8")
	ws.PushStreamServerInfo(&eventModel)
	return nil
}
