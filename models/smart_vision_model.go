package models

type SmartVisionModel struct {
	Id      string `json:"id" redis:"id"`
	Brand   string `json:"brand" redis:"brand"`
	Name    string `json:"name" redis:"name"`
	Address string `json:"address" redis:"address"`

	CreatedAt string `json:"created_at" redis:"created_at"`

	SelectedListJson string `json:"selected_list_json" redis:"selected_list_json"`
	ZonesList        string `json:"zones_list" redis:"zones_list"`
	MasksList        string `json:"masks_list" redis:"masks_list"`
	StartTime        string `json:"start_time" redis:"start_time"`
	EndTime          string `json:"end_time" redis:"end_time"`
}
