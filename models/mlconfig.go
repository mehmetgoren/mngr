package models

type MlConfig struct {
	DeviceConfig struct {
		DeviceName     string `json:"device_name"`
		DeviceServices []int  `json:"device_services"`
		DeviceType     int    `json:"device_type"`
	} `json:"device_config"`
	Handler struct {
		ReadServiceOverlay  bool   `json:"read_service_overlay"`
		SaveImageExtension  string `json:"save_image_extension"`
		SaveImageFolderPath string `json:"save_image_folder_path"`
		ShowImageCaption    bool   `json:"show_image_caption"`
		ShowImageFullscreen bool   `json:"show_image_fullscreen"`
		ShowImageWaitKey    int    `json:"show_image_wait_key"`
	} `json:"handler"`
	Heartbeat struct {
		Interval int `json:"interval"`
	} `json:"heartbeat"`
	Jetson struct {
		ModelName string  `json:"model_name"`
		Threshold float64 `json:"threshold"`
		WhiteList []int   `json:"white_list"`
	} `json:"jetson"`
	OnceDetector struct {
		ImagehashThreshold int     `json:"imagehash_threshold"`
		PsnrThreshold      float64 `json:"psnr_threshold"`
		SsimThreshold      float64 `json:"ssim_threshold"`
	} `json:"once_detector"`
	SourceReader struct {
		BufferSize      int  `json:"buffer_size"`
		Fps             int  `json:"fps"`
		KillStarterProc bool `json:"kill_starter_proc"`
		MaxRetry        int  `json:"max_retry"`
		MaxRetryIn      int  `json:"max_retry_in"`
		ReaderType      int  `json:"reader_type"`
	} `json:"source_reader"`
	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"redis"`
	Torch struct {
		ModelName         string  `json:"model_name"`
		ModelNameSpecific string  `json:"model_name_specific"`
		Threshold         float64 `json:"threshold"`
		WhiteList         []int   `json:"white_list"`
	} `json:"torch"`
}
