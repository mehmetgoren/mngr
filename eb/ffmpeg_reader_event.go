package eb

import (
	"github.com/go-redis/redis/v8"
	"mngr/models"
	"mngr/utils"
)

type FFmpegReaderResponseEvent struct {
	models.FFmpegReaderModel
	Pusher utils.WsPusher `json:"-"`
}

func (e *FFmpegReaderResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, e)
	e.Pusher.Push(e)
	return nil
}
