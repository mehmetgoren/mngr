package eb

import (
	"github.com/go-redis/redis/v8"
	"mngr/utils"
)

// EventHandler needs desperately generics
type EventHandler interface {
	Handle(event *redis.Message) error
	SetPusher(pusher utils.WsPusher)
}
