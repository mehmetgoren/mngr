package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/utils"
)

type StopStreamingRequestEvent struct {
	Id string `json:"id"`
}

func (s StopStreamingRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s *StopStreamingRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "stop_streaming_request"}
	err := eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a streaming event: " + err.Error())
		return err
	}

	return nil
}

type StopStreamingResponseEvent struct {
	Id     string         `json:"id"`
	Pusher utils.WsPusher `json:"-"`
}

func (s *StopStreamingResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, s)
	s.Pusher.Push(s)
	return nil
}
