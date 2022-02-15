package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/utils"
)

const (
	NONE               = 0
	TAKE_SCREENSHOT    = 1
	GENERATE_THUMBNAIL = 2
)

type EditorRequestEvent struct {
	Id          string `json:"id"`
	Brand       string `json:"brand"`
	Name        string `json:"name"`
	RtspAddress string `json:"rtsp_address"`
	EventType   int    `json:"event_type"`
}

func (e EditorRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}
func (e *EditorRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "editor_request"}
	err := eventBusPub.Publish(e)
	if err != nil {
		log.Println("An error occurred while publishing a editor event: " + err.Error())
		return err
	}

	return nil
}

type EditorResponseEvent struct {
	EditorRequestEvent
	ImageBase64 string         `json:"image_base64"`
	Pusher      utils.WsPusher `json:"-"`
}

func (e *EditorResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, e)
	push := func() {
		e.Pusher.Push(e)
	}
	switch e.EventType {
	case TAKE_SCREENSHOT:
		push()
		break
	case GENERATE_THUMBNAIL:
		push()
		break
	}
	return nil
}
