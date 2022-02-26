package models

type StreamModel struct {
	Id          string `json:"id" redis:"id"`
	Brand       string `json:"brand" redis:"brand"`
	Name        string `json:"name" redis:"name"`
	RtspAddress string `json:"rtsp_address" redis:"rtsp_address"`

	Pid       int    `json:"pid" redis:"pid"`
	CreatedAt string `json:"created_at" redis:"created_at"`
	Args      string `json:"args" redis:"args"`

	StreamType              int    `json:"stream_type" redis:"stream_type"`
	RtmpServerInitialized   bool   `json:"rtmp_server_initialized" redis:"rtmp_server_initialized"`
	RtmpServerType          int    `json:"rtmp_server_type" redis:"rtmp_server_type"`
	FlvPlayerConnectionType int    `json:"flv_player_connection_type" redis:"flv_player_connection_type"`
	RtmpImageName           string `json:"rtmp_image_name" redis:"rtmp_image_name"`
	RtmpContainerName       string `json:"rtmp_container_name" redis:"rtmp_container_name"`
	RtmpAddress             string `json:"rtmp_address" redis:"rtmp_address"`
	RtmpFlvAddress          string `json:"rtmp_flv_address" redis:"rtmp_flv_address"`
	RtmpContainerPorts      string `json:"rtmp_container_ports" redis:"rtmp_container_ports"`
	RtmpContainerCommands   string `json:"rtmp_container_commands" redis:"rtmp_container_commands"`

	DirectReadFrameRate int `json:"direct_read_frame_rate" redis:"direct_read_frame_rate"`
	DirectReadWidth     int `json:"direct_read_width" redis:"direct_read_width"`
	DirectReadHeight    int `json:"direct_read_height" redis:"direct_read_height"`

	JpegEnabled   bool `json:"jpeg_enabled" redis:"jpeg_enabled"`
	JpegFrameRate int  `json:"jpeg_frame_rate" redis:"jpeg_frame_rate"`

	Record               bool   `json:"record" redis:"record"`
	RecordDuration       int    `json:"record_duration" redis:"record_duration"`
	RecordFlvPid         int    `json:"record_flv_pid" redis:"record_flv_pid"`
	RecordFlvArgs        string `json:"record_flv_args" redis:"record_flv_args"`
	RecordFlvFailedCount int    `json:"record_flv_failed_count" redis:"record_flv_failed_count"`

	UseDiskImageReaderService bool `json:"use_disk_image_reader_service" redis:"use_disk_image_reader_service"`
	Reader                    bool `json:"reader" redis:"reader"`
	ReaderFrameRate           int  `json:"reader_frame_rate" redis:"reader_frame_rate"`
	ReaderWidth               int  `json:"reader_width" redis:"reader_width"`
	ReaderHeight              int  `json:"reader_height" redis:"reader_height"`
	ReaderPid                 int  `json:"reader_pid" redis:"reader_pid"`
	ReaderFailedCount         int  `json:"reader_failed_count" redis:"reader_failed_count"`

	HlsOutputPath          string `json:"hls_output_path" redis:"hls_output_path"`
	ReadJpegOutputPath     string `json:"read_jpeg_output_path" redis:"read_jpeg_output_path"`
	RecordOutputFolderPath string `json:"record_output_folder_path" redis:"record_output_folder_path"`
}
