package eb

import (
	"mngr/utils"
)

var connPubSub = utils.CreateRedisConnection(utils.EVENTBUS)

func SubscribeStartRecordingEvents(pusherRecording utils.WsPusher) {
	recordingEventBus := EventBus{Connection: connPubSub, Channel: "start_recording_response"}
	recordingEvent := StartRecordingEvent{Pusher: pusherRecording}
	go recordingEventBus.Subscribe(&recordingEvent)
}

func SubscribeStopRecordingEvents(pusherRecording utils.WsPusher) {
	recordingEventBus := EventBus{Connection: connPubSub, Channel: "stop_recording_response"}
	recordingEvent := StopRecordingEvent{Pusher: pusherRecording}
	go recordingEventBus.Subscribe(&recordingEvent)
}

func SubscribeStartStreamingEvents(pusherStreaming utils.WsPusher) {
	streamingEventBusSub := EventBus{Connection: connPubSub, Channel: "start_streaming_response"}
	streamingEventSub := StartStreamingEvent{Pusher: pusherStreaming}
	go streamingEventBusSub.Subscribe(&streamingEventSub)
}

func SubscribeStopStreamingEvents(pusherStreaming utils.WsPusher) {
	streamingEventBusSub := EventBus{Connection: connPubSub, Channel: "stop_streaming_response"}
	streamingEventSub := StopStreamingEvent{Pusher: pusherStreaming}
	go streamingEventBusSub.Subscribe(&streamingEventSub)
}
