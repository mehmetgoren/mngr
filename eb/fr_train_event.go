package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

type FaceTrainRequestEvent struct {
	Rb *reps.RepoBucket `json:"-"`
}

func (f *FaceTrainRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: f.Rb.PubSubConnection, Channel: "fr_train_request"}

	js, err := json.Marshal(f)
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

type FaceTrainResponseEvent struct {
	Result bool           `json:"result"`
	Pusher utils.WsPusher `json:"-"`
}

func (f *FaceTrainResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, f)
	f.Pusher.Push(f)
	return nil
}
