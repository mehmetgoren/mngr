package models

type AiModuleModel struct {
	Name                   string  `json:"name" redis:"name"`
	Description            string  `json:"description" redis:"description"`
	Enabled                bool    `json:"enabled" redis:"enabled"`
	ApiUrl                 string  `json:"api_url" redis:"api_url"`
	Threshold              float32 `json:"threshold" redis:"threshold"`
	LabelField             string  `json:"label_field" redis:"label_field"`
	MotionDetectionEnabled bool    `json:"motion_detection_enabled" redis:"motion_detection_enabled"`
	PersistenceEnabled     bool    `json:"persistence_enabled" redis:"persistence_enabled"`
	NotificationEnabled    bool    `json:"notification_enabled" redis:"notification_enabled"`
}
