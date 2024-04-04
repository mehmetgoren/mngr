package ws

import (
	"mngr/eb"
)

const (
	User         = 0
	StartStream  = 1
	StopStream   = 2
	Editor       = 3
	FFmpegReader = 4
	VideoMerge   = 5
	FaceTrain    = 6
	Probe        = 7
	Notifier     = 8
)

type UserEvent struct {
}

func (UserEvent) GetOp() int {
	return User
}
func (UserEvent) GetChannelName(string) string {
	return ""
}
func (UserEvent) CreateEventHandler() eb.EventHandler {
	return nil
}

// StartStreamEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type StartStreamEvent struct {
}

func (StartStreamEvent) GetOp() int {
	return StartStream
}

func (StartStreamEvent) GetChannelName(keyExtended string) string {
	return "start_stream_response"
}

func (StartStreamEvent) CreateEventHandler() eb.EventHandler {
	return &eb.StartStreamResponseEvent{}
}

// StopStreamEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type StopStreamEvent struct {
}

func (StopStreamEvent) GetOp() int {
	return StopStream
}

func (StopStreamEvent) GetChannelName(string) string {
	return "stop_stream_response"
}

func (StopStreamEvent) CreateEventHandler() eb.EventHandler {
	return &eb.StopStreamResponseEvent{}
}

// EditorEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type EditorEvent struct {
}

func (EditorEvent) GetOp() int {
	return Editor
}

func (EditorEvent) GetChannelName(string) string {
	return "editor_response"
}

func (EditorEvent) CreateEventHandler() eb.EventHandler {
	return &eb.EditorResponseEvent{}
}

// FFmpegReaderEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type FFmpegReaderEvent struct {
}

func (FFmpegReaderEvent) GetOp() int {
	return FFmpegReader
}

func (FFmpegReaderEvent) GetChannelName(sourceId string) string {
	return "ffrs" + sourceId
}

func (FFmpegReaderEvent) CreateEventHandler() eb.EventHandler {
	return &eb.FFmpegReaderResponseEvent{}
}

// VideoMergeEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type VideoMergeEvent struct {
}

func (VideoMergeEvent) GetOp() int {
	return VideoMerge
}

func (VideoMergeEvent) GetChannelName(string) string {
	return "vfm_response"
}

func (VideoMergeEvent) CreateEventHandler() eb.EventHandler {
	return &eb.VfmResponseEvent{}
}

// FaceTrainEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type FaceTrainEvent struct {
}

func (FaceTrainEvent) GetOp() int {
	return FaceTrain
}

func (FaceTrainEvent) GetChannelName(string) string {
	return "face_train_response"
}

func (FaceTrainEvent) CreateEventHandler() eb.EventHandler {
	return &eb.FaceTrainResponseEvent{}
}

// ProbeEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type ProbeEvent struct {
}

func (ProbeEvent) GetOp() int {
	return Probe
}

func (ProbeEvent) GetChannelName(string) string {
	return "probe_response"
}

func (ProbeEvent) CreateEventHandler() eb.EventHandler {
	return &eb.ProbeResponseEvent{}
}

// NotifierEvent ////////////////////////////////////////////////////////////////////////////////////////////////////////
type NotifierEvent struct {
}

func (NotifierEvent) GetOp() int {
	return Notifier
}

func (NotifierEvent) GetChannelName(string) string {
	return "notifier"
}

func (NotifierEvent) CreateEventHandler() eb.EventHandler {
	return &eb.NotifierResponseEvent{}
}
