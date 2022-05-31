package models

type User struct {
	Id          string `json:"id" redis:"id"`
	Username    string `json:"username" redis:"username"`
	Password    string `json:"password" redis:"password"`
	LastLoginAt string `json:"last_login_at" redis:"last_login_at"`
	Token       string `json:"token" redis:"token"`
	Email       string `json:"email" redis:"email"`
}

type RegisterUserViewModel struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RePassword string `json:"re_password" validate:"required"`
	Email      string `json:"email" validate:"email"`
}

type LoginUserViewModel struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
