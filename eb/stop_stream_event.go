package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

type StopStreamRequestEvent struct {
	Id string           `json:"id"`
	Rb *reps.RepoBucket `json:"-"`
}

func (s StopStreamRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s *StopStreamRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: s.Rb.PubSubConnection, Channel: "stop_stream_request"}
	err := eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a stream event: " + err.Error())
		return err
	}

	return nil
}

type StopStreamResponseEvent struct {
	Id     string         `json:"id"`
	Pusher utils.WsPusher `json:"-"`
}

func (s *StopStreamResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, s)
	s.Pusher.Push(s)
	return nil
}
