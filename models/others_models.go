package models

type RtspTemplateModel struct {
	Id              string `json:"id" redis:"id"`
	Name            string `json:"name" redis:"name"`
	Description     string `json:"description" redis:"description"`
	Brand           string `json:"brand" redis:"brand"`
	DefaultUser     string `json:"default_user" redis:"default_user"`
	DefaultPassword string `json:"default_password" redis:"default_password"`
	DefaultPort     string `json:"default_port" redis:"default_port"`
	Route           string `json:"route" redis:"route"`
	Templates       string `json:"templates" redis:"templates"`
}

type FailedStreamModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	RtmpContainerFailedCount int    `json:"rtmp_container_failed_count" redis:"rtmp_container_failed_count"`
	RtmpFeederFailedCount    int    `json:"rtmp_feeder_failed_count" redis:"rtmp_feeder_failed_count"`
	HlsFailedCount           int    `json:"hls_failed_count" redis:"hls_failed_count"`
	FFmpegReaderFailedCount  int    `json:"ffmpeg_reader_failed_count" redis:"ffmpeg_reader_failed_count"`
	RecordFailedCount        int    `json:"record_failed_count" redis:"record_failed_count"`
	SnapshotFailedCount      int    `json:"snapshot_failed_count" redis:"snapshot_failed_count"`
	RecordStuckProcessCount  int    `json:"record_stuck_process_count" redis:"record_stuck_process_count"`
	SourceStateConflictCount int    `json:"source_state_conflict_count" redis:"source_state_conflict_count"`
	LastCheckAt              string `json:"last_check_at" redis:"last_check_at"`
}

type RecStuckModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	RecordSegmentInterval int    `json:"record_segment_interval" redis:"record_segment_interval"`
	RecordOutputDir       string `json:"record_output_dir" redis:"record_output_dir"`
	FileExt               string `json:"file_ext" redis:"file_ext"`
	LastModifiedFile      string `json:"last_modified_file" redis:"last_modified_file"`
	LastModifiedSize      int    `json:"last_modified_size" redis:"last_modified_size"`
	FailedCount           int    `json:"failed_count" redis:"failed_count"`
	FailedModifiedFile    string `json:"failed_modified_file" redis:"failed_modified_file"`
	LastCheckAt           string `json:"last_check_at" redis:"last_check_at"`
}

type VariousInfos struct {
	RtmpPortCounter      int      `json:"rtmp_port_counter"`
	RtmpContainerZombies []string `json:"rtmp_container_zombies"`
	FFmpegProcessZombies []int    `json:"ffmpeg_process_zombies"`
}
