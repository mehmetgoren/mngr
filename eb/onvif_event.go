package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/models"
	"mngr/reps"
	"mngr/utils"
)

// OnvifRequestEvent it is just a simple proxy
type OnvifRequestEvent struct {
	models.OnvifEvent
	Rb *reps.RepoBucket `json:"-"`
}

func (o OnvifRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *OnvifRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: o.Rb.PubSubConnection, Channel: "onvif_request"}

	js, err := json.Marshal(o)
	if err != nil {
		return err
	}

	err = eventBusPub.Publish(string(js))
	if err != nil {
		log.Println("An error occurred while publishing a restart stream event: " + err.Error())
		return err
	}

	return nil
}

type OnvifResponseEvent struct {
	models.OnvifEvent
	Pusher utils.WsPusher `json:"-"`
}

func (o *OnvifResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, o)
	o.Pusher.Push(o)
	return nil
}
