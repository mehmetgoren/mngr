package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

type ProbeRequestEvent struct {
	Address string           `json:"address"`
	Rb      *reps.RepoBucket `json:"-"`
}

func (p ProbeRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProbeRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: p.Rb.PubSubConnection, Channel: "probe_request"}
	err := eventBusPub.Publish(p)
	if err != nil {
		log.Println("An error occurred while publishing a probe event: " + err.Error())
		return err
	}

	return nil
}

type ProbeResponseEvent struct {
	Address   string         `json:"address"`
	ResultB64 string         `json:"result_b64"`
	Pusher    utils.WsPusher `json:"-"`
}

func (p *ProbeResponseEvent) Handle(event *redis.Message) error {
	err := utils.DeserializeJson(event.Payload, p)
	if err != nil {
		log.Println("An error occurred while handling a probe event: " + err.Error())
		return err
	}
	err = p.Pusher.Push(p)
	if err != nil {
		log.Println("An error occurred while handling a probe event: " + err.Error())
		return err
	}
	return nil
}

func (p *ProbeResponseEvent) SetPusher(pusher utils.WsPusher) {
	p.Pusher = pusher
}
