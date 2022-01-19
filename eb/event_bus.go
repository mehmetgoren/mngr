package eb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/utils"
)

var ConnPubSub = utils.CreateRedisConnection(utils.EVENTBUS)

type EventBus struct {
	Connection *redis.Client
	Channel    string
}

func (eb *EventBus) Publish(event interface{}) error {
	return eb.Connection.Publish(context.Background(), eb.Channel, event).Err()
}

func (eb *EventBus) Subscribe(handler EventHandler) error {
	pong, err := eb.Connection.Ping(context.Background()).Result()
	if err != nil {
		log.Println("ping has been failed, exiting now...")
		panic(err)
		return err
	}

	log.Println("ping: " + pong + " for " + eb.Channel)
	log.Println("redis pubsub is listening for " + eb.Channel)

	channel := eb.Channel
	subscribe := eb.Connection.Subscribe(context.Background(), channel)
	subscriptions := subscribe.ChannelWithSubscriptions(context.Background(), 1)
	for {
		select {
		case sub := <-subscriptions:
			var message, isRedisMessage = sub.(*redis.Message)
			if !isRedisMessage {
				continue
			}
			go func() {
				err := handler.Handle(message)
				if err != nil {
					log.Println("an error has been occurred while handling the event: ", err)
				}
			}()
		}
	}
}
