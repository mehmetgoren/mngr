package models

type OdModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	CreatedAt string `json:"created_at" redis:"created_at"`

	ThresholdList string `json:"threshold_list" redis:"threshold_list"`
	SelectedList  string `json:"selected_list" redis:"selected_list"`
	ZonesList     string `json:"zones_list" redis:"zones_list"`
	MasksList     string `json:"masks_list" redis:"masks_list"`
	StartTime     string `json:"start_time" redis:"start_time"`
	EndTime       string `json:"end_time" redis:"end_time"`
}
