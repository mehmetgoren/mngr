package models

import "time"

type User struct {
	Id                 string `json:"id" redis:"id"`
	Username           string `json:"username" redis:"username"`
	Password           string `json:"password" redis:"password"`
	LastLoginAt        string `json:"last_login_at" redis:"last_login_at"`
	Token              string `json:"token" redis:"token"`
	Email              string `json:"email" redis:"email"`
	Ip                 string `json:"ip" redis:"ip"`
	Uag                string `json:"uag" redis:"uag"`
	Location           string `json:"location" redis:"location"`
	DataCenterLocation string `json:"data_center_location" redis:"data_center_location"`
}

type RegisterUserViewModel struct {
	Username           string `json:"username" validate:"required"`
	Password           string `json:"password" validate:"required"`
	RePassword         string `json:"re_password" validate:"required"`
	Email              string `json:"email" validate:"email"`
	Ip                 string `json:"ip" redis:"ip"`
	Uag                string `json:"uag" redis:"uag"`
	Location           string `json:"location" redis:"location"`
	DataCenterLocation string `json:"data_center_location" redis:"data_center_location"`
}

type LoginUserViewModel struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserSession struct {
	*User
	LastVisitAt time.Time
}
