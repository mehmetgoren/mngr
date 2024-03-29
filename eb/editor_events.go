package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

const (
	NONE               = 0
	TAKE_SCREENSHOT    = 1
	GENERATE_THUMBNAIL = 2
	MASK_SCREENSHOT    = 3
)

type EditorRequestEvent struct {
	Id        string           `json:"id"`
	Brand     string           `json:"brand"`
	Name      string           `json:"name"`
	Address   string           `json:"address"`
	EventType int              `json:"event_type"`
	Rb        *reps.RepoBucket `json:"-"`
}

func (e EditorRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}
func (e *EditorRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: e.Rb.PubSubConnection, Channel: "editor_request"}
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
	case MASK_SCREENSHOT:
		push()
		break
	}
	return nil
}

func (e *EditorResponseEvent) SetPusher(pusher utils.WsPusher) {
	e.Pusher = pusher
}
