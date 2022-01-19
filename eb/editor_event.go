package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/utils"
)

const (
	NONE               = 0
	TAKE_SCREENSHOT    = 1
	GENERATE_THUMBNAIL = 2
)

type EditorEvent struct {
	Source       models.Source  `json:"source"`
	EventType    int            `json:"event_type"`
	ResponseJson string         `json:"response_json"`
	Pusher       utils.WsPusher `json:"-"`
}

func (e EditorEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}
func (e EditorEvent) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}

func (e *EditorEvent) Publish() error {
	eventBusPub := EventBus{Connection: ConnPubSub, Channel: "editor_request"}
	err := eventBusPub.Publish(e)
	if err != nil {
		log.Println("An error occurred while publishing a editor event: " + err.Error())
		return err
	}

	return nil
}

func (e *EditorEvent) Handle(event *redis.Message) error {
	var eventModel EditorEvent
	utils.DeserializeJson(event.Payload, &eventModel)
	switch eventModel.EventType {
	case TAKE_SCREENSHOT:
		responseModel := models.EditorImageResponseModel{}
		responseModel.EventType = eventModel.EventType
		responseModel.Source = eventModel.Source
		utils.DeserializeJson(eventModel.ResponseJson, &responseModel)
		e.Pusher.Push(&responseModel)
		break
	}
	return nil
}
