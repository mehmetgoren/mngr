package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type StopStreamingEvent struct {
	models.Source
	Pusher utils.WsPusher `json:"-"`
}

func (s StopStreamingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s StopStreamingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *StopStreamingEvent) Publish() error {
	eventBusPub := EventBus{Connection: connPubSub, Channel: "stop_streaming_request"}
	err := eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a streaming event: " + err.Error())
		return err
	}

	return nil
}

func (s *StopStreamingEvent) Handle(event *redis.Message) error {
	var eventModel StopStreamingEvent
	utils.DeserializeJson(event.Payload, &eventModel)
	s.Pusher.Push(&eventModel)
	return nil
}
