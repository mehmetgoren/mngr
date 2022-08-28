package eb

import (
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/utils"
)

type NotifierResponseEvent struct {
	Base64Object string         `json:"base_64_object"`
	Type         int            `json:"type"`
	Pusher       utils.WsPusher `json:"-"`
}

func (n *NotifierResponseEvent) Handle(event *redis.Message) error {
	err := utils.DeserializeJson(event.Payload, n)
	if err != nil {
		log.Println("An error occurred while handling a notifier event: " + err.Error())
		return err
	}
	err = n.Pusher.Push(n)
	if err != nil {
		log.Println("An error occurred while handling a notifier event: " + err.Error())
		return err
	}
	return nil
}

func (n *NotifierResponseEvent) SetPusher(pusher utils.WsPusher) {
	n.Pusher = pusher
}
