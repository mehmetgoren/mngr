package eb

import (
	"encoding/json"
	"log"
	"mngr/models"
	"mngr/utils"
)

type RestartStreamRequestEvent struct {
	models.SourceModel
}

func (r RestartStreamRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
func (r *RestartStreamRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "restart_stream_request"}
	err := eventBusPub.Publish(r)
	if err != nil {
		log.Println("An error occurred while publishing a restart stream event: " + err.Error())
		return err
	}

	return nil
}

// restart does not need to be subscribed to since it is only called by the restart_stream_request which is just a proxy.
