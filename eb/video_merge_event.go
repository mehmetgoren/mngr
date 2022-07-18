package eb

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"mngr/reps"
	"mngr/utils"
)

type VfmRequestEvent struct {
	SourceId string           `json:"source_id"`
	DateStr  string           `json:"date_str"`
	Rb       *reps.RepoBucket `json:"-"`
}

func (v VfmRequestEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VfmRequestEvent) Publish() error {
	eventBusPub := EventBus{PubSubConnection: v.Rb.PubSubConnection, Channel: "vfm_request"}
	err := eventBusPub.Publish(v)
	if err != nil {
		log.Println("An error occurred while publishing a stream event: " + err.Error())
		return err
	}

	return nil
}

type VfmResponseEvent struct {
	SourceId             string         `json:"source_id"`
	OutputFileName       string         `json:"output_file_name"`
	MergedVideoFilenames []string       `json:"merged_video_filenames"`
	Pusher               utils.WsPusher `json:"-"`
}

func (v *VfmResponseEvent) Handle(event *redis.Message) error {
	utils.DeserializeJson(event.Payload, v)
	v.Pusher.Push(v)
	return nil
}

func (v *VfmResponseEvent) SetPusher(pusher utils.WsPusher) {
	v.Pusher = pusher
}
