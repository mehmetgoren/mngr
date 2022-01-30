package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
	"strings"
)

type StartStreamingRequestEvent struct {
	models.SourceModel
}

func (s StartStreamingRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s *StartStreamingRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "start_streaming_request"}
	err := eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a streaming event: " + err.Error())
		return err
	}

	return nil
}

type StartStreamingResponseEvent struct {
	models.StreamingModel
	Pusher utils.WsPusher `json:"-"`
}

func (s *StartStreamingResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, s)
	//from full path to web server relative path
	config, _ := utils.ConfigRep.GetConfig()
	s.HlsOutputPath = strings.Replace(s.HlsOutputPath, config.Path.Streaming, "", -1)
	s.Pusher.Push(s)
	return nil
}
