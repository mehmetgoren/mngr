package models

type Source struct {
	TypeName    string `json:"type_name"`
	Name        string `json:"name"`
	Brand       string `json:"brand"`
	RtspAddress string `json:"rtsp_address"`
	Id          string `json:"id"`
}
