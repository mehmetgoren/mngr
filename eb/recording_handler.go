package eb

import (
	"github.com/go-redis/redis/v8"
	"mngr/utils"
	"mngr/ws"
)

func (r *RecordingEvent) Handle(event *redis.Message) error {
	var eventModel RecordingEvent
	utils.DeserializeJson(event.Payload, &eventModel)

	eventModel.OutputFile = eventModel.Id
	eventModel.Duration = 3
	ws.PushRecordServerInfo(&eventModel)
	return nil
}
