package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type StopRecordingEvent struct {
	models.Source
	Pusher utils.WsPusher `json:"-"`
}

func (r StopRecordingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
func (r StopRecordingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *StopRecordingEvent) Publish() error {
	eventBusPub := EventBus{Connection: connPubSub, Channel: "stop_recording_request"}
	err := eventBusPub.Publish(r)
	if err != nil {
		log.Println("An error occurred while publishing a playback event: " + err.Error())
		return err
	}

	return nil
}

func (r *StopRecordingEvent) Handle(event *redis.Message) error {
	var eventModel StartRecordingEvent
	utils.DeserializeJson(event.Payload, &eventModel)
	r.Pusher.Push(&eventModel)
	return nil
}
