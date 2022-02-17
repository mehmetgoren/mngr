package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type StartStreamRequestEvent struct {
	models.SourceModel
}

func (s StartStreamRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s *StartStreamRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "start_stream_request"}
	err := eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a stream event: " + err.Error())
		return err
	}

	return nil
}

type StartStreamResponseEvent struct {
	models.StreamModel
	Pusher utils.WsPusher `json:"-"`
}

func (s *StartStreamResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, s)
	//from full path to web server relative path
	config, _ := utils.ConfigRep.GetConfig()
	utils.SetHlsPath(config, &s.StreamModel)
	s.Pusher.Push(s)
	return nil
}
