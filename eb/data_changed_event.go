package eb

import (
	"encoding/json"
	"log"
	"mngr/reps"
)

const (
	SAVE   = 0
	DELETE = 1
)

type ModelChanged struct {
	SourceId string `json:"source_id" redis:"source_id"`
}

type DataChangedEvent struct {
	ModelName  string           `json:"model_name" redis:"model_name"`
	ParamsJson string           `json:"params_json" redis:"params_json"`
	Op         int              `json:"op" redis:"op"`
	Rb         *reps.RepoBucket `json:"-"`
}

func (d DataChangedEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *DataChangedEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: d.Rb.PubSubConnection, Channel: "data_changed"}
	err := eventBusPub.Publish(d)
	if err != nil {
		log.Println("An error occurred while publishing a data changed event: " + err.Error())
		return err
	}

	return nil
}
