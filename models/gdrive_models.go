package models

type GdriveViewModel struct {
	AuthCode        string `json:"auth_code"`
	CredentialsJson string `json:"credentials_json"`
	Enabled         bool   `json:"enabled"`
	TokenJson       string `json:"token_json"`
	URL             string `json:"url"`
}
