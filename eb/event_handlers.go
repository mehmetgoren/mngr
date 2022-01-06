package eb

import (
	"github.com/go-redis/redis/v8"
	"mngr/utils"
	"mngr/ws"
	"path/filepath"
)

// EventHandler needs desperately generics
type EventHandler interface {
	Handle(event *redis.Message) error
}

func (h *StreamingEvent) Handle(event *redis.Message) error {
	var eventModel StreamingEvent
	utils.DeserializeJson(event.Payload, &eventModel)

	eventModel.FolderPath = filepath.Join(eventModel.Id, "stream.m3u8")
	ws.PushStreamServerInfo(&eventModel)
	return nil
}
