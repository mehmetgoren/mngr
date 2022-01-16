package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

type RecordingEvent struct {
	models.Source
	Duration   int            `json:"duration"`
	OutputFile string         `json:"output_file"`
	Pusher     utils.WsPusher `json:"-"`
}

func (r RecordingEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
func (r RecordingEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (r *RecordingEvent) Publish() error {
	folderPath, err := utils.CreateDirIfNotExist(utils.PlaybackFolderPath + "/" + r.Source.Id)
	if err != nil {
		log.Println("An error occurred while creating a playback folder: " + err.Error())
		return err
	}

	folderPathFull, _ := utils.GetExecutablePath()
	//todo: move it to config
	r.Duration = 3
	r.OutputFile = folderPathFull + "/" + folderPath

	eventBusPub := EventBus{Connection: connPubSub, Channel: "start_recording_request"}
	err = eventBusPub.Publish(r)
	if err != nil {
		log.Println("An error occurred while publishing a playback event: " + err.Error())
		return err
	}

	return nil
}

func (r *RecordingEvent) Handle(event *redis.Message) error {
	var eventModel RecordingEvent
	utils.DeserializeJson(event.Payload, &eventModel)

	eventModel.OutputFile = eventModel.Id
	eventModel.Duration = 3
	r.Pusher.Push(&eventModel)
	//ws.PushRecordServerInfo(&eventModel)
	return nil
}
