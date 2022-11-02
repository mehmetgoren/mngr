package models

import (
	"encoding/json"
)

type SourceModel struct {
	// FFmpegModel section starts
	Id            string `json:"id" redis:"id"`
	Address       string `json:"address" redis:"address"`
	RtspTransport int    `json:"rtsp_transport" redis:"rtsp_transport"`

	AnalyzationDuration int    `json:"analyzation_duration" redis:"analyzation_duration"`
	ProbeSize           int    `json:"probe_size" redis:"probe_size"`
	InputFrameRate      int    `json:"input_frame_rate" redis:"input_frame_rate"`
	UseCameraTimestamp  bool   `json:"use_camera_timestamp" redis:"use_camera_timestamp"`
	UseHwaccel          bool   `json:"use_hwaccel" redis:"use_hwaccel"`
	HwaccelEngine       int    `json:"hwaccel_engine" redis:"hwaccel_engine"`
	VideoDecoder        int    `json:"video_decoder" redis:"video_decoder"`
	HwaccelDevice       string `json:"hwaccel_device" redis:"hwaccel_device"`

	StreamType            int    `json:"stream_type" redis:"stream_type"`
	RtmpAddress           string `json:"rtmp_address" redis:"rtmp_address"`
	StreamVideoCodec      int    `json:"stream_video_codec" redis:"stream_video_codec"`
	Preset                int    `json:"preset" redis:"preset"`
	HlsTime               int    `json:"hls_time" redis:"hls_time"`
	HlsListSize           int    `json:"hls_list_size" redis:"hls_list_size"`
	StreamQuality         int    `json:"stream_quality" redis:"stream_quality"`
	StreamFrameRate       int    `json:"stream_frame_rate" redis:"stream_frame_rate"`
	StreamWidth           int    `json:"stream_width" redis:"stream_width"`
	StreamHeight          int    `json:"stream_height" redis:"stream_height"`
	StreamRotate          int    `json:"stream_rotate" redis:"stream_rotate"`
	StreamAudioCodec      int    `json:"stream_audio_codec" redis:"stream_audio_codec"`
	StreamAudioChannel    int    `json:"stream_audio_channel" redis:"stream_audio_channel"`
	StreamAudioQuality    int    `json:"stream_audio_quality" redis:"stream_audio_quality"`
	StreamAudioSampleRate int    `json:"stream_audio_sample_rate" redis:"stream_audio_sample_rate"`
	StreamAudioVolume     int    `json:"stream_audio_volume" redis:"stream_audio_volume"`

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
	// FFmpegModel section ends

	// SourceModel section starts
	Brand       string `json:"brand" redis:"brand"`
	Name        string `json:"name" redis:"name"`
	Description string `json:"description" redis:"description"`

	MacAddress string `json:"mac_address" redis:"mac_address"`
	IpAddress  string `json:"ip_address" redis:"ip_address"`

	Enabled        bool `json:"enabled" redis:"enabled"`
	State          int  `json:"state" redis:"state"`
	RtmpServerType int  `json:"rtmp_server_type" redis:"rtmp_server_type"`

	SnapshotEnabled      bool    `json:"snapshot_enabled" redis:"snapshot_enabled"`
	SnapshotType         int     `json:"snapshot_type" redis:"snapshot_type"`
	SnapshotFrameRate    int     `json:"snapshot_frame_rate" redis:"snapshot_frame_rate"`
	SnapshotWidth        int     `json:"snapshot_width" redis:"snapshot_width"`
	SnapshotHeight       int     `json:"snapshot_height" redis:"snapshot_height"`
	MdType               int     `json:"md_type" redis:"md_type"`
	MdOpencvThreshold    int     `json:"md_opencv_threshold" redis:"md_opencv_threshold"`
	MdContourAreaLimit   int     `json:"md_contour_area_limit" redis:"md_contour_area_limit"`
	MdImagehashThreshold int     `json:"md_imagehash_threshold" redis:"md_imagehash_threshold"`
	MdPsnrThreshold      float32 `json:"md_psnr_threshold" redis:"md_psnr_threshold"`

	FfmpegReaderFrameRate int `json:"ffmpeg_reader_frame_rate" redis:"ffmpeg_reader_frame_rate"`
	FfmpegReaderWidth     int `json:"ffmpeg_reader_width" redis:"ffmpeg_reader_width"`
	FfmpegReaderHeight    int `json:"ffmpeg_reader_height" redis:"ffmpeg_reader_height"`

	RecordEnabled bool `json:"record_enabled" redis:"record_enabled"`
	AiClipEnabled bool `json:"ai_clip_enabled" redis:"ai_clip_enabled"`

	FlvPlayerType  int  `json:"flv_player_type" redis:"flv_player_type"`
	BoosterEnabled bool `json:"booster_enabled" redis:"booster_enabled"`

	BlackScreenCheckEnabled bool   `json:"black_screen_check_enabled" redis:"black_screen_check_enabled"`
	CreatedAt               string `json:"created_at" redis:"created_at"`
	// SourceModel section ends
}

func (s SourceModel) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}
