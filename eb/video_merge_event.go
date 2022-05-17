package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

type VideMergeRequestEvent struct {
	Id      string           `json:"id"`
	DateStr string           `json:"date_str"`
	Rb      *reps.RepoBucket `json:"-"`
}

func (v VideMergeRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VideMergeRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: v.Rb.PubSubConnection, Channel: "vfm_request"}
	err := eventBusPub.Publish(v)
	if err != nil {
		log.Println("An error occurred while publishing a stream event: " + err.Error())
		return err
	}

	return nil
}

type VideMergeResponseEvent struct {
	Id     string         `json:"id"`
	Result bool           `json:"result"`
	Pusher utils.WsPusher `json:"-"`
}

func (v *VideMergeResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, v)
	v.Pusher.Push(v)
	return nil
}
