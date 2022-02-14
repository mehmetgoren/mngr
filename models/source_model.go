package models

import (
	"encoding/json"
)

type SourceModel struct {
	Id          string `json:"id" redis:"id"`
	Brand       string `json:"brand" redis:"brand"`
	Name        string `json:"name" redis:"name"`
	RtspAddress string `json:"rtsp_address" redis:"rtsp_address"`
	Description string `json:"description" redis:"description"`

	Enabled       bool `json:"enabled" redis:"enabled"`
	Record        bool `json:"record" redis:"record"`
	InputType     int  `json:"input_type" redis:"input_type"`
	RtspTransport int  `json:"rtsp_transport" redis:"rtsp_transport"`

	AnalyzationDuration int    `json:"analyzation_duration" redis:"analyzation_duration"`
	ProbeSize           int    `json:"probe_size" redis:"probe_size"`
	InputFrameRate      int    `json:"input_frame_rate" redis:"input_frame_rate"`
	UseCameraTimestamp  bool   `json:"use_camera_timestamp" redis:"use_camera_timestamp"`
	UseHwaccel          bool   `json:"use_hwaccel" redis:"use_hwaccel"`
	HwaccelEngine       int    `json:"hwaccel_engine" redis:"hwaccel_engine"`
	VideoDecoder        int    `json:"video_decoder" redis:"video_decoder"`
	HwaccelDevice       string `json:"hwaccel_device" redis:"hwaccel_device"`

	StreamType              int    `json:"stream_type" redis:"stream_type"`
	RtmpServerType          int    `json:"rtmp_server_type" redis:"rtmp_server_type"`
	FlvPlayerConnectionType int    `json:"flv_player_connection_type" redis:"flv_player_connection_type"`
	RtmpServerAddress       string `json:"rtmp_server_address" redis:"rtmp_server_address"`
	NeedReloadInterval      int    `json:"need_reload_interval" redis:"need_reload_interval"`
	DirectReadFrameRate     int    `json:"direct_read_frame_rate" redis:"direct_read_frame_rate"`
	DirectReadWidth         int    `json:"direct_read_width" redis:"direct_read_width"`
	DirectReadHeight        int    `json:"direct_read_height" redis:"direct_read_height"`
	StreamVideoCodec        int    `json:"stream_video_codec" redis:"stream_video_codec"`
	HlsTime                 int    `json:"hls_time" redis:"hls_time"`
	HlsListSize             int    `json:"hls_list_size" redis:"hls_list_size"`
	HlsPreset               int    `json:"hls_preset" redis:"hls_preset"`
	StreamQuality           int    `json:"stream_quality" redis:"stream_quality"`
	StreamFrameRate         int    `json:"stream_frame_rate" redis:"stream_frame_rate"`
	StreamWidth             int    `json:"stream_width" redis:"stream_width"`
	StreamHeight            int    `json:"stream_height" redis:"stream_height"`
	StreamRotate            int    `json:"stream_rotate" redis:"stream_rotate"`
	StreamAudioCodec        int    `json:"stream_audio_codec" redis:"stream_audio_codec"`
	StreamAudioChannel      int    `json:"stream_audio_channel" redis:"stream_audio_channel"`
	StreamAudioQuality      int    `json:"stream_audio_quality" redis:"stream_audio_quality"`
	StreamAudioSampleRate   int    `json:"stream_audio_sample_rate" redis:"stream_audio_sample_rate"`
	StreamAudioVolume       int    `json:"stream_audio_volume" redis:"stream_audio_volume"`

	JpegEnabled               bool `json:"jpeg_enabled" redis:"jpeg_enabled"`
	JpegFrameRate             int  `json:"jpeg_frame_rate" redis:"jpeg_frame_rate"`
	JpegUseVsync              bool `json:"jpeg_use_vsync" redis:"jpeg_use_vsync"`
	JpegQuality               int  `json:"jpeg_quality" redis:"jpeg_quality"`
	JpegWidth                 int  `json:"jpeg_width" redis:"jpeg_width"`
	JpegHeight                int  `json:"jpeg_height" redis:"jpeg_height"`
	UseDiskImageReaderService bool `json:"use_disk_image_reader_service" redis:"use_disk_image_reader_service"`

	RecordFileType        int `json:"record_file_type" redis:"record_file_type"`
	RecordVideoCodec      int `json:"record_video_codec" redis:"record_video_codec"`
	RecordQuality         int `json:"record_quality" redis:"record_quality"`
	RecordPreset          int `json:"record_preset" redis:"record_preset"`
	RecordFrameRate       int `json:"record_frame_rate" redis:"record_frame_rate"`
	RecordWidth           int `json:"record_width" redis:"record_width"`
	RecordHeight          int `json:"record_height" redis:"record_height"`
	RecordSegmentInterval int `json:"record_segment_interval" redis:"record_segment_interval"`
	RecordRotate          int `json:"record_rotate" redis:"record_rotate"`
	RecordAudioCodec      int `json:"record_audio_codec" redis:"record_audio_codec"`
	RecordAudioChannel    int `json:"record_audio_channel" redis:"record_audio_channel"`
	RecordAudioQuality    int `json:"record_audio_quality" redis:"record_audio_quality"`
	RecordAudioSampleRate int `json:"record_audio_sample_rate" redis:"record_audio_sample_rate"`
	RecordAudioVolume     int `json:"record_audio_volume" redis:"record_audio_volume"`

	LogLevel int `json:"log_level" redis:"log_level"`
}

func (s SourceModel) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
