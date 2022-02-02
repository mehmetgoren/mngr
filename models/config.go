package models

type Config struct {
	DeviceConfig struct {
		DeviceName     string `json:"device_name"`
		DeviceType     int    `json:"device_type"`
		DeviceServices []int  `json:"device_services"`
	} `json:"device"`
	HeartbeatConfig struct {
		Interval int `json:"interval"`
	} `json:"heartbeat"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Jetson struct {
		ModelName string  `json:"model_name"`
		Threshold float64 `json:"threshold"`
		WhiteList []int   `json:"white_list"`
	} `json:"jetson"`
	Torch struct {
		ModelName         string  `json:"model_name"`
		ModelNameSpecific string  `json:"model_name_specific"`
		Threshold         float64 `json:"threshold"`
		WhiteList         []int   `json:"white_list"`
	} `json:"torch"`
	OnceDetector struct {
		ImagehashThreshold int     `json:"imagehash_threshold"`
		PsnrThreshold      float64 `json:"psnr_threshold"`
		SsimThreshold      float64 `json:"ssim_threshold"`
	} `json:"once_detector"`
	Handler struct {
		SaveImageFolderPath string `json:"save_image_folder_path"`
		SaveImageExtension  string `json:"save_image_extension"`
		ShowImageWaitKey    int    `json:"show_image_wait_key"`
		ShowImageCaption    bool   `json:"show_image_caption"`
		ShowImageFullscreen bool   `json:"show_image_fullscreen"`
		ReadServiceOverlay  bool   `json:"read_service_overlay"`
	} `json:"handler"`
	SourceReader struct {
		Fps        int `json:"fps"`
		BufferSize int `json:"buffer_size"`
		MaxRetry   int `json:"max_retry"`
		MaxRetryIn int `json:"max_retry_in"`
		// todo: remove it.
		KillStarterProc bool `json:"kill_starter_proc"`
	} `json:"source_reader"`
	Path struct {
		Streaming string `json:"streaming"`
		Recording string `json:"recording"`
		Reading   string `json:"reading"`
	} `json:"path"`
	FFmpeg struct {
		UseDoubleQuotesForPath bool `json:"use_double_quotes_for_path"`
		MaxOperationRetryCount int  `json:"max_operation_retry_count"`
	} `json:"ffmpeg"`
}
