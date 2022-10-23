package models

type TelegramBot struct {
	Token string `json:"token" redis:"token"`
	URL   string `json:"url" redis:"url"`
}

type TelegramUser struct {
	ID int64 `json:"id" redis:"id"`

	FirstName    string `json:"first_name" redis:"first_name"`
	LastName     string `json:"last_name" redis:"last_name"`
	Username     string `json:"username" redis:"username"`
	LanguageCode string `json:"language_code" redis:"language_code"`
	IsBot        bool   `json:"is_bot" redis:"is_bot"`

	// Returns only in getMe
	CanJoinGroups   bool `json:"can_join_groups" redis:"can_join_groups"`
	CanReadMessages bool `json:"can_read_all_group_messages" redis:"can_read_all_group_messages"`
	SupportsInline  bool `json:"supports_inline_queries" redis:"supports_inline_queries"`
}

type TelegramViewModel struct {
	Enabled bool            `json:"enabled"`
	Bot     *TelegramBot    `json:"bot"`
	Users   []*TelegramUser `json:"users"`
}
