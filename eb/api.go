package eb

import (
	"mngr/utils"
)

var connPubSub = utils.CreateRedisConnection(utils.EVENTBUS)

func SubscribeStreamingEvents(pusherStreaming utils.WsPusher) {
	streamingEventBusSub := EventBus{Connection: connPubSub, Channel: "start_streaming_response"}
	streamingEventSub := StreamingEvent{Pusher: pusherStreaming}
	go streamingEventBusSub.Subscribe(&streamingEventSub)
}

func SubscribeRecordingEvents(pusherRecording utils.WsPusher) {
	recordingEventBus := EventBus{Connection: connPubSub, Channel: "start_recording_response"}
	recordingEvent := RecordingEvent{Pusher: pusherRecording}
	go recordingEventBus.Subscribe(&recordingEvent)
}
