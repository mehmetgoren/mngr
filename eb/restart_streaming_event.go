package eb

import (
	"encoding/json"
	"log"
	"mngr/models"
	"mngr/utils"
)

type RestartStreamingRequestEvent struct {
	models.SourceModel
}

func (r RestartStreamingRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
func (r *RestartStreamingRequestEvent) Publish() error {
	eventBusPub := EventBus{Connection: utils.ConnPubSub, Channel: "restart_streaming_request"}
	err := eventBusPub.Publish(r)
	if err != nil {
		log.Println("An error occurred while publishing a restart streaming event: " + err.Error())
		return err
	}

	return nil
}

// restart does not need to be subscribed to since it is only called by the restart_streaming_request which is just a proxy.
