package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
	"path/filepath"
)

type StartStreamingEvent struct {
	models.Source
	OutputFile string         `json:"output_file"`
	Pusher     utils.WsPusher `json:"-"`
}

func (s StartStreamingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
func (s StartStreamingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s *StartStreamingEvent) Publish() error {
	folderPath, err := utils.CreateDirIfNotExist(utils.LiveFolderPath + "/" + s.Source.Id)
	if err != nil {
		log.Println("An error occurred while creating a live stream folder: " + err.Error())
		return err
	}

	folderPathFull, _ := utils.GetExecutablePath()
	s.OutputFile = folderPathFull + "/" + folderPath + "/stream.m3u8"
	//create file
	//os.Create(eventPub.OutputFile)
	eventBusPub := EventBus{Connection: ConnPubSub, Channel: "start_streaming_request"}
	err = eventBusPub.Publish(s)
	if err != nil {
		log.Println("An error occurred while publishing a streaming event: " + err.Error())
		return err
	}

	return nil
}

func (s *StartStreamingEvent) Handle(event *redis.Message) error {
	var eventModel StartStreamingEvent
	utils.DeserializeJson(event.Payload, &eventModel)

	eventModel.OutputFile = filepath.Join(eventModel.Id, "stream.m3u8")
	s.Pusher.Push(&eventModel)
	//ws.PushStreamServerInfo(&eventModel)
	return nil
}
